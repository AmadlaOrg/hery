package git

import (
	"github.com/stretchr/testify/mock"
)

// MockGit is the mock implementation of IGit
type MockGit struct {
	mock.Mock
}

func (m *MockGit) FetchRepo(url, dest string) error {
	args := m.Called(url, dest)
	return args.Error(0)
}

func (m *MockGit) CommitHeadHash(repoPath string) (string, error) {
	args := m.Called(repoPath)
	return args.String(0), args.Error(1)
}

func (m *MockGit) CheckoutTag(repoPath, tagName string) error {
	args := m.Called(repoPath, tagName)
	return args.Error(0)
}

// FIXME:

// Test FetchRepo
/*func TestFetchRepo(t *testing.T) {
	mockGit := new(MockGit)
	testURL := "https://github.com/git-fixtures/basic.git"
	testDest := "/tmp/repo"

	mockGit.On("FetchRepo", testURL, testDest).Return(nil)

	sgit := &SGit{}

	err := sgit.FetchRepo(testURL, testDest)

	// Assert that the method was called with correct params
	mockGit.AssertCalled(t, "FetchRepo", testURL, testDest)
	require.NoError(t, err)
}*/

// Test FetchRepo with Error
/*func TestFetchRepo_Error(t *testing.T) {
	mockGit := new(MockGit)
	testURL := "https://github.com/git-fixtures/basic.git"
	testDest := "/tmp/repo"

	mockGit.On("FetchRepo", testURL, testDest).Return(errors.New("failed to clone repo"))

	sgit := &SGit{}

	err := sgit.FetchRepo(testURL, testDest)

	// Assert that the method was called with correct params
	mockGit.AssertCalled(t, "FetchRepo", testURL, testDest)
	require.Error(t, err)
	require.EqualError(t, err, "failed to clone repo")
}

// Test CommitHeadHash
func TestCommitHeadHash(t *testing.T) {
	mockGit := new(MockGit)
	testRepoPath := "/tmp/repo"
	expectedHash := "abc123"

	mockGit.On("CommitHeadHash", testRepoPath).Return(expectedHash, nil)

	sgit := &SGit{}

	hash, err := sgit.CommitHeadHash(testRepoPath)

	// Assert that the method was called with correct params
	mockGit.AssertCalled(t, "CommitHeadHash", testRepoPath)
	require.NoError(t, err)
	require.Equal(t, expectedHash, hash)
}

// Test CommitHeadHash with Error
func TestCommitHeadHash_Error(t *testing.T) {
	mockGit := new(MockGit)
	testRepoPath := "/tmp/repo"

	mockGit.On("CommitHeadHash", testRepoPath).Return("", errors.New("failed to get commit hash"))

	sgit := &SGit{}

	hash, err := sgit.CommitHeadHash(testRepoPath)

	// Assert that the method was called with correct params
	mockGit.AssertCalled(t, "CommitHeadHash", testRepoPath)
	require.Error(t, err)
	require.Equal(t, "", hash)
	require.EqualError(t, err, "failed to get commit hash")
}

// Test CheckoutTag
func TestCheckoutTag(t *testing.T) {
	mockGit := new(MockGit)
	testRepoPath := "/tmp/repo"
	testTag := "v1.0.0"

	mockGit.On("CheckoutTag", testRepoPath, testTag).Return(nil)

	sgit := &SGit{}

	err := sgit.CheckoutTag(testRepoPath, testTag)

	// Assert that the method was called with correct params
	mockGit.AssertCalled(t, "CheckoutTag", testRepoPath, testTag)
	require.NoError(t, err)
}

// Test CheckoutTag with Error
func TestCheckoutTag_Error(t *testing.T) {
	mockGit := new(MockGit)
	testRepoPath := "/tmp/repo"
	testTag := "v1.0.0"

	mockGit.On("CheckoutTag", testRepoPath, testTag).Return(errors.New("failed to checkout tag"))

	sgit := &SGit{}

	err := sgit.CheckoutTag(testRepoPath, testTag)

	// Assert that the method was called with correct params
	mockGit.AssertCalled(t, "CheckoutTag", testRepoPath, testTag)
	require.Error(t, err)
	require.EqualError(t, err, "failed to checkout tag")
}*/
