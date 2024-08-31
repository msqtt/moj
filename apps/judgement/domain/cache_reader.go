package domain

import (
	"context"
	"errors"
	"io"
	"log"
	"moj/domain/judgement"
	"moj/domain/policy"
	"moj/domain/question"
	"moj/judgement/etc"
	"net/http"
	"net/url"

	"github.com/valkey-io/valkey-go"
)

type CacheCaseReader struct {
	client valkey.Client
	conf   *etc.Config
}

// Read implements policy.CaseFileService.
func (m *CacheCaseReader) Read(ctx context.Context, filePath string) (fileContent string, err error) {
	cmd := m.client.B().Get().Key(filePath).Build()
	cacheStr, err := m.client.Do(ctx, cmd).ToString()
	if err == nil {
		fileContent = cacheStr
		return
	}
	url, err := url.JoinPath(m.conf.OssPrefixURL, filePath)
	if err != nil {
		return
	}
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		err = errors.Join(err, errors.New("failed to get file from oss"), err)
		return
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Join(err, errors.New("failed to read file from response"), err)
		return
	}

	fileContent = string(bytes)

	// set cache
	go func() {
		cmd = m.client.B().Set().Key(filePath).Value(fileContent).
			ExSeconds(int64(m.conf.CaseCacheExpireTime.Seconds())).Build()
		m.client.Do(ctx, cmd)
	}()

	return
}

// ReadAllCaseFile implements policy.CaseFileService.
func (m *CacheCaseReader) ReadAllCaseFile(ctx context.Context, caseFiles []question.Case) ([]judgement.Case, error) {
	ret := make([]judgement.Case, len(caseFiles))
	var errs []error
	for id, f := range caseFiles {
		inputContent, err := m.Read(ctx, f.InputFilePath)
		if err != nil {
			errs = append(errs, err)
		}
		outputContent, err := m.Read(ctx, f.OutputFilePath)
		if err != nil {
			errs = append(errs, err)
		}
		ret[id] = judgement.Case{
			Number:       f.Number,
			InputContent: inputContent,
			OuputContent: outputContent,
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return ret, nil
}

func NewCacheCaseReader(conf *etc.Config) policy.CaseFileService {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{conf.ValkeyHostAddr}})
	if err != nil {
		log.Fatal(client)
	}
	return &CacheCaseReader{
		client: client,
		conf:   conf,
	}
}
