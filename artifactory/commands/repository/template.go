package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/jfrog/jfrog-cli-core/v2/artifactory/commands/utils"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/utils/ioutils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type RepoTemplateCommand struct {
	path string
}

const (
	// Strings for prompt questions
	SelectConfigKeyMsg = "Select the next configuration key" + ioutils.PressTabMsg

	TemplateType = "templateType"
	Create       = "create"
	Update       = "update"

	MandatoryUrl = "mandatoryUrl"

	// Common repository configuration JSON keys
	Key             = "key"
	Rclass          = "rclass"
	PackageType     = "packageType"
	Description     = "description"
	Notes           = "notes"
	IncludePatterns = "includesPattern"
	ExcludePatterns = "excludesPattern"
	RepoLayoutRef   = "repoLayoutRef"
	ProjectKey      = "projectKey"
	Environment     = "environment"

	// Mutual local and remote repository configuration JSON keys
	HandleReleases               = "handleReleases"
	HandleSnapshots              = "handleSnapshots"
	MaxUniqueSnapshots           = "maxUniqueSnapshots"
	SuppressPomConsistencyChecks = "suppressPomConsistencyChecks"
	BlackedOut                   = "blackedOut"
	XrayIndex                    = "xrayIndex"
	PropertySets                 = "propertySets"
	DownloadRedirect             = "downloadRedirect"
	BlockPushingSchema1          = "blockPushingSchema1"

	// Mutual local and virtual repository configuration JSON keys
	DebianTrivialLayout             = "debianTrivialLayout"
	OptionalIndexCompressionFormats = "optionalIndexCompressionFormats"
	PrimaryKeyPairRef               = "primaryKeyPairRef"

	// Mutual remote and virtual repository configuration JSON keys
	ExternalDependenciesEnabled  = "externalDependenciesEnabled"
	ExternalDependenciesPatterns = "externalDependenciesPatterns"

	// Unique local repository configuration JSON keys
	ChecksumPolicyType       = "checksumPolicyType"
	MaxUniqueTags            = "maxUniqueTags"
	SnapshotVersionBehavior  = "snapshotVersionBehavior"
	ArchiveBrowsingEnabled   = "archiveBrowsingEnabled"
	CalculateYumMetadata     = "calculateYumMetadata"
	YumRootDepth             = "yumRootDepth"
	DockerApiVersion         = "dockerApiVersion"
	EnableFileListsIndexing  = "enableFileListsIndexing"
	ForceNugetAuthentication = "forceNugetAuthentication"

	// Unique remote repository configuration JSON keys
	Url                               = "url"
	Username                          = "username"
	Password                          = "password"
	Proxy                             = "proxy"
	RemoteRepoChecksumPolicyType      = "remoteRepoChecksumPolicyType"
	HardFail                          = "hardFail"
	Offline                           = "offline"
	StoreArtifactsLocally             = "storeArtifactsLocally"
	SocketTimeoutMillis               = "socketTimeoutMillis"
	LocalAddress                      = "localAddress"
	RetrievalCachePeriodSecs          = "retrievalCachePeriodSecs"
	FailedRetrievalCachePeriodSecs    = "failedRetrievalCachePeriodSecs"
	MissedRetrievalCachePeriodSecs    = "missedRetrievalCachePeriodSecs"
	UnusedArtifactsCleanupEnabled     = "unusedArtifactsCleanupEnabled"
	UnusedArtifactsCleanupPeriodHours = "unusedArtifactsCleanupPeriodHours"
	AssumedOfflinePeriodSecs          = "assumedOfflinePeriodSecs"
	FetchJarsEagerly                  = "fetchJarsEagerly"
	FetchSourcesEagerly               = "fetchSourcesEagerly"
	RejectInvalidJars                 = "rejectInvalidJars"
	ShareConfiguration                = "shareConfiguration"
	SynchronizeProperties             = "synchronizeProperties"
	BlockMismatchingMimeTypes         = "blockMismatchingMimeTypes"
	AllowAnyHostAuth                  = "allowAnyHostAuth"
	EnableCookieManagement            = "enableCookieManagement"
	BowerRegistryUrl                  = "bowerRegistryUrl"
	ComposerRegistryUrl               = "composerRegistryUrl"
	PyPIRegistryUrl                   = "pyPIRegistryUrl"
	VcsType                           = "vcsType"
	VcsGitProvider                    = "vcsGitProvider"
	VcsGitDownloadUrl                 = "vcsGitDownloadUrl"
	BypassHeadRequests                = "bypassHeadRequests"
	ClientTlsCertificate              = "clientTlsCertificate"
	FeedContextPath                   = "feedContextPath"
	DownloadContextPath               = "downloadContextPath"
	V3FeedUrl                         = "v3FeedUrl"
	ContentSynchronisation            = "contentSynchronisation"
	ListRemoteFolderItems             = "listRemoteFolderItems"
	EnableTokenAuthentication         = "enableTokenAuthentication"
	PodsSpecsRepoUrl                  = "podsSpecsRepoUrl"

	// Unique virtual repository configuration JSON keys
	Repositories                                  = "repositories"
	ArtifactoryRequestsCanRetrieveRemoteArtifacts = "artifactoryRequestsCanRetrieveRemoteArtifacts"
	KeyPair                                       = "keyPair"
	PomRepositoryReferencesCleanupPolicy          = "pomRepositoryReferencesCleanupPolicy"
	DefaultDeploymentRepo                         = "defaultDeploymentRepo"
	ForceMavenAuthentication                      = "forceMavenAuthentication"
	ExternalDependenciesRemoteRepo                = "externalDependenciesRemoteRepo"

	// rclasses
	Local     = "local"
	Remote    = "remote"
	Virtual   = "virtual"
	Federated = "federated"

	// PackageTypes
	Generic   = "generic"
	Maven     = "maven"
	Gradle    = "gradle"
	Ivy       = "ivy"
	Sbt       = "sbt"
	Helm      = "helm"
	Cocoapods = "cocoapods"
	Opkg      = "opkg"
	Rpm       = "rpm"
	Nuget     = "nuget"
	Cran      = "cran"
	Gems      = "gems"
	Npm       = "npm"
	Bower     = "bower"
	Debian    = "debian"
	Composer  = "composer"
	Pypi      = "pypi"
	Docker    = "docker"
	Vagrant   = "vagrant"
	Gitlfs    = "gitlfs"
	Go        = "go"
	Yum       = "yum"
	Conan     = "conan"
	Chef      = "chef"
	Puppet    = "puppet"
	Vcs       = "vcs"
	Alpine    = "alpine"
	Conda     = "conda"
	P2        = "p2"

	// Repo layout Refs
	BowerDefaultRepoLayout    = "bower-default"
	buildDefaultRepoLayout    = "build-default"
	ComposerDefaultRepoLayout = "composer-default"
	ConanDefaultRepoLayout    = "conan-default"
	GoDefaultRepoLayout       = "go-default"
	GradleDefaultRepoLayout   = "maven-2-default"
	IvyDefaultRepoLayout      = "ivy-default"
	Maven1DefaultRepoLayout   = "maven-1-default"
	Maven2DefaultRepoLayout   = "maven-2-default"
	NpmDefaultRepoLayout      = "npm-default"
	NugetDefaultRepoLayout    = "nuget-default"
	puppetDefaultRepoLayout   = "puppet-default"
	SbtDefaultRepoLayout      = "sbt-default"
	SimpleDefaultRepoLayout   = "simple-default"
	VcsDefaultRepoLayout      = "vcs-default"

	// Checksum Policies
	ClientChecksumPolicy           = "client-checksums"
	ServerGeneratedChecksumsPolicy = "server-generated-checksums"

	// Snapshot version behaviors
	UniqueBehavior    = "unique"
	NonUniqueBehavior = "non-unique"
	DeployerBehavior  = "deployer"

	// Optional index compression formats
	Bz2Compression  = "bz2"
	LzmaCompression = "lzma"
	XzCompression   = "xz"

	// Docker api versions
	DockerApiV1 = "V1"
	DockerApiV2 = "V2"

	// Remote repo checksum policies
	GenerateIfAbsentPolicy  = "generate-if-absent"
	FailPolicy              = "fail"
	IgnoreAndGeneratePolicy = "ignore-and-generate"
	PassThruPolicy          = "pass-thru"

	// Vcs Types
	Git = "GIT"

	// Vcs git provider
	GithubVcsProvider      = "GITHUB"
	BitbucketVcsProvider   = "BITBUCKET"
	OldstashVcsProvider    = "OLDSTASH"
	StashVcsProvider       = "STASH"
	ArtifactoryVcsProvider = "ARTIFACTORY"
	CustomVcsProvider      = "CUSTOM"

	// POM repository references cleanup policies
	DiscardActiveRefrencePolicy = "discard_active_reference"
	DiscardAnyReferencePolicy   = "discard_any_reference"
	NothingPolicy               = "nothing"
)

var optionalSuggestsMap = map[string]prompt.Suggest{
	ioutils.SaveAndExit:               {Text: ioutils.SaveAndExit},
	Description:                       {Text: Description},
	Notes:                             {Text: Notes},
	IncludePatterns:                   {Text: IncludePatterns},
	ExcludePatterns:                   {Text: ExcludePatterns},
	RepoLayoutRef:                     {Text: RepoLayoutRef},
	ProjectKey:                        {Text: ProjectKey},
	Environment:                       {Text: Environment},
	HandleReleases:                    {Text: HandleReleases},
	HandleSnapshots:                   {Text: HandleSnapshots},
	MaxUniqueSnapshots:                {Text: MaxUniqueSnapshots},
	SuppressPomConsistencyChecks:      {Text: SuppressPomConsistencyChecks},
	BlackedOut:                        {Text: BlackedOut},
	DownloadRedirect:                  {Text: DownloadRedirect},
	BlockPushingSchema1:               {Text: BlockPushingSchema1},
	DebianTrivialLayout:               {Text: DebianTrivialLayout},
	ExternalDependenciesEnabled:       {Text: ExternalDependenciesEnabled},
	ExternalDependenciesPatterns:      {Text: ExternalDependenciesPatterns},
	ChecksumPolicyType:                {Text: ChecksumPolicyType},
	MaxUniqueTags:                     {Text: MaxUniqueTags},
	SnapshotVersionBehavior:           {Text: SnapshotVersionBehavior},
	XrayIndex:                         {Text: XrayIndex},
	PropertySets:                      {Text: PropertySets},
	ArchiveBrowsingEnabled:            {Text: ArchiveBrowsingEnabled},
	CalculateYumMetadata:              {Text: CalculateYumMetadata},
	YumRootDepth:                      {Text: YumRootDepth},
	DockerApiVersion:                  {Text: DockerApiVersion},
	EnableFileListsIndexing:           {Text: EnableFileListsIndexing},
	OptionalIndexCompressionFormats:   {Text: OptionalIndexCompressionFormats},
	Url:                               {Text: Url},
	Username:                          {Text: Username},
	Password:                          {Text: Password},
	Proxy:                             {Text: Proxy},
	PrimaryKeyPairRef:                 {Text: PrimaryKeyPairRef},
	RemoteRepoChecksumPolicyType:      {Text: RemoteRepoChecksumPolicyType},
	HardFail:                          {Text: HardFail},
	Offline:                           {Text: Offline},
	StoreArtifactsLocally:             {Text: StoreArtifactsLocally},
	SocketTimeoutMillis:               {Text: SocketTimeoutMillis},
	LocalAddress:                      {Text: LocalAddress},
	RetrievalCachePeriodSecs:          {Text: RetrievalCachePeriodSecs},
	FailedRetrievalCachePeriodSecs:    {Text: FailedRetrievalCachePeriodSecs},
	MissedRetrievalCachePeriodSecs:    {Text: MissedRetrievalCachePeriodSecs},
	UnusedArtifactsCleanupEnabled:     {Text: UnusedArtifactsCleanupEnabled},
	UnusedArtifactsCleanupPeriodHours: {Text: UnusedArtifactsCleanupPeriodHours},
	AssumedOfflinePeriodSecs:          {Text: AssumedOfflinePeriodSecs},
	FetchJarsEagerly:                  {Text: FetchJarsEagerly},
	FetchSourcesEagerly:               {Text: FetchSourcesEagerly},
	RejectInvalidJars:                 {Text: RejectInvalidJars},
	ShareConfiguration:                {Text: ShareConfiguration},
	SynchronizeProperties:             {Text: SynchronizeProperties},
	BlockMismatchingMimeTypes:         {Text: BlockMismatchingMimeTypes},
	AllowAnyHostAuth:                  {Text: AllowAnyHostAuth},
	EnableCookieManagement:            {Text: EnableCookieManagement},
	BowerRegistryUrl:                  {Text: BowerRegistryUrl},
	ComposerRegistryUrl:               {Text: ComposerRegistryUrl},
	PyPIRegistryUrl:                   {Text: PyPIRegistryUrl},
	VcsType:                           {Text: VcsType},
	VcsGitProvider:                    {Text: VcsGitProvider},
	VcsGitDownloadUrl:                 {Text: VcsGitDownloadUrl},
	BypassHeadRequests:                {Text: BypassHeadRequests},
	ClientTlsCertificate:              {Text: ClientTlsCertificate},
	FeedContextPath:                   {Text: FeedContextPath},
	DownloadContextPath:               {Text: DownloadContextPath},
	V3FeedUrl:                         {Text: V3FeedUrl},
	ContentSynchronisation:            {Text: ContentSynchronisation},
	ListRemoteFolderItems:             {Text: ListRemoteFolderItems},
	PodsSpecsRepoUrl:                  {Text: PodsSpecsRepoUrl},
	EnableTokenAuthentication:         {Text: EnableTokenAuthentication},
	Repositories:                      {Text: Repositories},
	ArtifactoryRequestsCanRetrieveRemoteArtifacts: {Text: ArtifactoryRequestsCanRetrieveRemoteArtifacts},
	KeyPair:                              {Text: KeyPair},
	PomRepositoryReferencesCleanupPolicy: {Text: PomRepositoryReferencesCleanupPolicy},
	DefaultDeploymentRepo:                {Text: DefaultDeploymentRepo},
	ForceMavenAuthentication:             {Text: ForceMavenAuthentication},
	ForceNugetAuthentication:             {Text: ForceNugetAuthentication},
	ExternalDependenciesRemoteRepo:       {Text: ExternalDependenciesRemoteRepo},
}

var baseLocalRepoConfKeys = []string{
	Description, Notes, IncludePatterns, ExcludePatterns, RepoLayoutRef, ProjectKey, Environment, BlackedOut, XrayIndex,
	PropertySets, ArchiveBrowsingEnabled, OptionalIndexCompressionFormats, DownloadRedirect, BlockPushingSchema1,
}

var mavenGradleLocalRepoConfKeys = []string{
	MaxUniqueSnapshots, HandleReleases, HandleSnapshots, SuppressPomConsistencyChecks, SnapshotVersionBehavior, ChecksumPolicyType,
}

var rpmLocalRepoConfKeys = []string{
	YumRootDepth, CalculateYumMetadata, EnableFileListsIndexing, PrimaryKeyPairRef,
}

var nugetLocalRepoConfKeys = []string{
	MaxUniqueSnapshots, ForceNugetAuthentication,
}

var debianLocalRepoConfKeys = []string{
	DebianTrivialLayout, PrimaryKeyPairRef,
}

var dockerLocalRepoConfKeys = []string{
	DockerApiVersion, MaxUniqueTags,
}

var baseRemoteRepoConfKeys = []string{
	Username, Password, Proxy, Description, Notes, IncludePatterns, ExcludePatterns, RepoLayoutRef, ProjectKey, Environment, HardFail, Offline,
	BlackedOut, XrayIndex, StoreArtifactsLocally, SocketTimeoutMillis, LocalAddress, RetrievalCachePeriodSecs, FailedRetrievalCachePeriodSecs,
	MissedRetrievalCachePeriodSecs, UnusedArtifactsCleanupEnabled, UnusedArtifactsCleanupPeriodHours, AssumedOfflinePeriodSecs,
	ShareConfiguration, SynchronizeProperties, BlockMismatchingMimeTypes, PropertySets, AllowAnyHostAuth, EnableCookieManagement,
	BypassHeadRequests, ClientTlsCertificate, DownloadRedirect, BlockPushingSchema1, ContentSynchronisation,
}

var mavenGradleRemoteRepoConfKeys = []string{
	FetchJarsEagerly, FetchSourcesEagerly, RemoteRepoChecksumPolicyType, HandleReleases, HandleSnapshots,
	SuppressPomConsistencyChecks, RejectInvalidJars,
}

var cocoapodsRemoteRepoConfKeys = []string{
	PodsSpecsRepoUrl,
}

var opkgRemoteRepoConfKeys = []string{
	ListRemoteFolderItems,
}

var rpmRemoteRepoConfKeys = []string{
	ListRemoteFolderItems,
}

var nugetRemoteRepoConfKeys = []string{
	FeedContextPath, DownloadContextPath, V3FeedUrl, ForceNugetAuthentication,
}

var gemsRemoteRepoConfKeys = []string{
	ListRemoteFolderItems,
}

var npmRemoteRepoConfKeys = []string{
	ListRemoteFolderItems,
}

var bowerRemoteRepoConfKeys = []string{
	BowerRegistryUrl,
}

var debianRemoteRepoConfKeys = []string{
	ListRemoteFolderItems,
}

var composerRemoteRepoConfKeys = []string{
	ComposerRegistryUrl,
}

var pypiRemoteRepoConfKeys = []string{
	PyPIRegistryUrl, ListRemoteFolderItems,
}

var dockerRemoteRepoConfKeys = []string{
	ExternalDependenciesEnabled, ExternalDependenciesPatterns, EnableTokenAuthentication,
}

var gitlfsRemoteRepoConfKeys = []string{
	ListRemoteFolderItems,
}

var vcsRemoteRepoConfKeys = []string{
	VcsGitProvider, VcsType, MaxUniqueSnapshots, VcsGitDownloadUrl, ListRemoteFolderItems,
}

var baseVirtualRepoConfKeys = []string{
	Repositories, Description, Notes, IncludePatterns, ExcludePatterns, RepoLayoutRef, ProjectKey, Environment, ArtifactoryRequestsCanRetrieveRemoteArtifacts,
	DefaultDeploymentRepo, OptionalIndexCompressionFormats, PrimaryKeyPairRef,
}

var mavenGradleVirtualRepoConfKeys = []string{
	ForceMavenAuthentication, PomRepositoryReferencesCleanupPolicy, KeyPair,
}

var nugetVirtualRepoConfKeys = []string{
	ForceNugetAuthentication,
}

var npmVirtualRepoConfKeys = []string{
	ExternalDependenciesEnabled, ExternalDependenciesPatterns, ExternalDependenciesRemoteRepo,
}

var bowerVirtualRepoConfKeys = []string{
	ExternalDependenciesEnabled, ExternalDependenciesPatterns, ExternalDependenciesRemoteRepo,
}

var debianVirtualRepoConfKeys = []string{
	DebianTrivialLayout,
}

var goVirtualRepoConfKeys = []string{
	ExternalDependenciesEnabled, ExternalDependenciesPatterns,
}

var commonPkgTypes = []string{
	Maven, Gradle, Ivy, Sbt, Helm, Rpm, Nuget, Cran, Gems, Npm, Bower, Debian, Pypi, Docker, Gitlfs, Go, Conan,
	Chef, Puppet, Alpine, Generic,
}

var localRepoAdditionalPkgTypes = []string{
	Cocoapods, Opkg, Composer, Vagrant, Yum,
}

var remoteRepoAdditionalPkgTypes = []string{
	Cocoapods, Opkg, Composer, Conda, P2, Vcs, Yum,
}

var virtualRepoAdditionalPkgTypes = []string{
	Conda, P2, Yum,
}

var federatedRepoAdditionalPkgTypes = []string{
	Vagrant, Opkg, Conda, Composer, Cocoapods,
}

var pkgTypeSuggestsMap = map[string]prompt.Suggest{
	Generic:   {Text: Generic},
	Maven:     {Text: Maven},
	Gradle:    {Text: Gradle},
	Ivy:       {Text: Ivy},
	Sbt:       {Text: Sbt},
	Helm:      {Text: Helm},
	Cocoapods: {Text: Cocoapods},
	Opkg:      {Text: Opkg},
	Rpm:       {Text: Rpm},
	Nuget:     {Text: Nuget},
	Cran:      {Text: Cran},
	Gems:      {Text: Gems},
	Npm:       {Text: Npm},
	Bower:     {Text: Bower},
	Debian:    {Text: Debian},
	Composer:  {Text: Composer},
	Pypi:      {Text: Pypi},
	Docker:    {Text: Docker},
	Vagrant:   {Text: Vagrant},
	Gitlfs:    {Text: Gitlfs},
	Go:        {Text: Go},
	Yum:       {Text: Yum},
	Conan:     {Text: Conan},
	Chef:      {Text: Chef},
	Puppet:    {Text: Puppet},
	Vcs:       {Text: Vcs},
	Conda:     {Text: Conda},
	P2:        {Text: P2},
	Alpine:    {Text: Alpine},
}

func NewRepoTemplateCommand() *RepoTemplateCommand {
	return &RepoTemplateCommand{}
}

func (rtc *RepoTemplateCommand) SetTemplatePath(path string) *RepoTemplateCommand {
	rtc.path = path
	return rtc
}

func (rtc *RepoTemplateCommand) ServerDetails() (*config.ServerDetails, error) {
	// Since it's a local command, usage won't be reported.
	return nil, nil
}

func (rtc *RepoTemplateCommand) Run() (err error) {
	err = utils.ValidateTemplatePath(rtc.path)
	if err != nil {
		return
	}
	repoTemplateQuestionnaire := &ioutils.InteractiveQuestionnaire{
		MandatoryQuestionsKeys: []string{TemplateType, Key, Rclass},
		QuestionsMap:           questionMap,
	}
	err = repoTemplateQuestionnaire.Perform()
	if err != nil {
		return err
	}
	resBytes, err := json.Marshal(repoTemplateQuestionnaire.AnswersMap)
	if err != nil {
		return errorutils.CheckError(err)
	}
	if err = os.WriteFile(rtc.path, resBytes, 0644); err != nil {
		return errorutils.CheckError(err)
	}
	log.Info(fmt.Sprintf("Repository configuration template successfully created at %s.", rtc.path))

	return nil
}

func (rtc *RepoTemplateCommand) CommandName() string {
	return "rt_repo_template"
}

func rclassCallback(iq *ioutils.InteractiveQuestionnaire, rclass string) (string, error) {
	var pkgTypes = commonPkgTypes
	switch rclass {
	case Remote:
		// For create template url is mandatory, for update we will allow url as an optional key
		if _, ok := iq.AnswersMap[TemplateType]; !ok {
			return "", errors.New("package type is missing in configuration map")
		}
		if iq.AnswersMap[TemplateType] == Create {
			_, err := iq.AskQuestion(iq.QuestionsMap[MandatoryUrl])
			if err != nil {
				return "", err
			}
		}
		pkgTypes = append(pkgTypes, remoteRepoAdditionalPkgTypes...)
	case Local:
		pkgTypes = append(pkgTypes, localRepoAdditionalPkgTypes...)
	case Virtual:
		pkgTypes = append(pkgTypes, virtualRepoAdditionalPkgTypes...)
	case Federated:
		pkgTypes = append(pkgTypes, federatedRepoAdditionalPkgTypes...)
	default:
		return "", errors.New("unsupported rclass")
	}
	// PackageType is also mandatory. Since the possible types depend on which rcalss was chosen, we ask the question here.
	var pkgTypeQuestion = ioutils.QuestionInfo{
		Options:      ioutils.GetSuggestsFromKeys(pkgTypes, pkgTypeSuggestsMap),
		Msg:          "",
		PromptPrefix: "Select the repository's package type" + ioutils.PressTabMsg,
		AllowVars:    false,
		Writer:       ioutils.WriteStringAnswer,
		MapKey:       PackageType,
		Callback:     pkgTypeCallback,
	}
	return iq.AskQuestion(pkgTypeQuestion)
}

func pkgTypeCallback(iq *ioutils.InteractiveQuestionnaire, pkgType string) (string, error) {
	// Each combination of (rclass,packageType) has its own optional configuration keys.
	// We set the questionnaire's optionalKeys suggests according to the selected combination.
	if _, ok := iq.AnswersMap[Rclass]; !ok {
		return "", errors.New("rclass is missing in configuration map")
	}
	switch iq.AnswersMap[Rclass] {
	case Local:
		iq.OptionalKeysSuggests = getLocalRepoConfKeys(pkgType)
	case Remote:
		// For update template we need to allow url as an optional key
		if _, ok := iq.AnswersMap[TemplateType]; !ok {
			return "", errors.New("package type is missing in configuration map")
		}
		iq.OptionalKeysSuggests = getRemoteRepoConfKeys(pkgType, fmt.Sprint(iq.AnswersMap[TemplateType]))
	case Virtual:
		iq.OptionalKeysSuggests = getVirtualRepoConfKeys(pkgType)
	case Federated:
		iq.OptionalKeysSuggests = getLocalRepoConfKeys(pkgType)
	default:
		return "", errors.New("unsupported rclass was configured")
	}
	// We don't need the templateType value in the final configuration
	delete(iq.AnswersMap, TemplateType)
	return "", nil
}

// Repo key must have a prefix of "<projectKey>-". This callback adds the prefix to the repo key if it is missing.
func projectKeyCallback(iq *ioutils.InteractiveQuestionnaire, projectKey string) (string, error) {
	if _, ok := iq.AnswersMap[Key]; !ok {
		return "", errorutils.CheckErrorf("repository key is missing in configuration map")
	}
	requiredProjectPrefix := projectKey + "-"
	currentRepoKey, ok := iq.AnswersMap[Key].(string)
	if !ok {
		return "", errorutils.CheckErrorf("template syntax error: the value for the repository key is not a string type")
	}

	if !strings.HasPrefix(currentRepoKey, requiredProjectPrefix) {
		newRepoKey := requiredProjectPrefix + currentRepoKey
		log.Info("Repository key should start with the projectKey followed by a dash. Modifying repo key to: '" + newRepoKey + "'.")
		iq.AnswersMap[Key] = newRepoKey
	}
	return "", nil
}

func getLocalRepoConfKeys(pkgType string) []prompt.Suggest {
	optionalKeys := []string{ioutils.SaveAndExit}
	optionalKeys = append(optionalKeys, baseLocalRepoConfKeys...)
	switch pkgType {
	case Maven, Gradle:
		optionalKeys = append(optionalKeys, mavenGradleLocalRepoConfKeys...)
	case Rpm:
		optionalKeys = append(optionalKeys, rpmLocalRepoConfKeys...)
	case Nuget:
		optionalKeys = append(optionalKeys, nugetLocalRepoConfKeys...)
	case Debian:
		optionalKeys = append(optionalKeys, debianLocalRepoConfKeys...)
	case Docker:
		optionalKeys = append(optionalKeys, dockerLocalRepoConfKeys...)
	}
	return ioutils.GetSuggestsFromKeys(optionalKeys, optionalSuggestsMap)
}

func getRemoteRepoConfKeys(pkgType, templateType string) []prompt.Suggest {
	optionalKeys := []string{ioutils.SaveAndExit}
	if templateType == Update {
		optionalKeys = append(optionalKeys, Url)
	}
	optionalKeys = append(optionalKeys, baseRemoteRepoConfKeys...)
	switch pkgType {
	case Maven, Gradle:
		optionalKeys = append(optionalKeys, mavenGradleRemoteRepoConfKeys...)
	case Cocoapods:
		optionalKeys = append(optionalKeys, cocoapodsRemoteRepoConfKeys...)
	case Opkg:
		optionalKeys = append(optionalKeys, opkgRemoteRepoConfKeys...)
	case Rpm:
		optionalKeys = append(optionalKeys, rpmRemoteRepoConfKeys...)
	case Nuget:
		optionalKeys = append(optionalKeys, nugetRemoteRepoConfKeys...)
	case Gems:
		optionalKeys = append(optionalKeys, gemsRemoteRepoConfKeys...)
	case Npm:
		optionalKeys = append(optionalKeys, npmRemoteRepoConfKeys...)
	case Bower:
		optionalKeys = append(optionalKeys, bowerRemoteRepoConfKeys...)
	case Debian:
		optionalKeys = append(optionalKeys, debianRemoteRepoConfKeys...)
	case Composer:
		optionalKeys = append(optionalKeys, composerRemoteRepoConfKeys...)
	case Pypi:
		optionalKeys = append(optionalKeys, pypiRemoteRepoConfKeys...)
	case Docker:
		optionalKeys = append(optionalKeys, dockerRemoteRepoConfKeys...)
	case Gitlfs:
		optionalKeys = append(optionalKeys, gitlfsRemoteRepoConfKeys...)
	case Vcs:
		optionalKeys = append(optionalKeys, vcsRemoteRepoConfKeys...)
	}
	return ioutils.GetSuggestsFromKeys(optionalKeys, optionalSuggestsMap)
}

func getVirtualRepoConfKeys(pkgType string) []prompt.Suggest {
	optionalKeys := []string{ioutils.SaveAndExit}
	optionalKeys = append(optionalKeys, baseVirtualRepoConfKeys...)
	switch pkgType {
	case Maven, Gradle:
		optionalKeys = append(optionalKeys, mavenGradleVirtualRepoConfKeys...)
	case Nuget:
		optionalKeys = append(optionalKeys, nugetVirtualRepoConfKeys...)
	case Npm:
		optionalKeys = append(optionalKeys, npmVirtualRepoConfKeys...)
	case Bower:
		optionalKeys = append(optionalKeys, bowerVirtualRepoConfKeys...)
	case Debian:
		optionalKeys = append(optionalKeys, debianVirtualRepoConfKeys...)
	case Go:
		optionalKeys = append(optionalKeys, goVirtualRepoConfKeys...)
	}
	return ioutils.GetSuggestsFromKeys(optionalKeys, optionalSuggestsMap)
}

func contentSynchronisationCallBack(iq *ioutils.InteractiveQuestionnaire, answer string) (value string, err error) {
	// contentSynchronisation has an object value with 4 bool fields.
	// We ask for the rest of the values and writes the values in comma separated list.
	if err != nil {
		return "", nil
	}
	answer += "," + ioutils.AskFromList("", "Insert the value for statistic.enable >", false, ioutils.GetBoolSuggests(), "")
	// cs.Statistics.Enabled, err = strconv.ParseBool(enabled)
	if err != nil {
		return "", nil
	}
	answer += "," + ioutils.AskFromList("", "Insert the value for properties.enable >", false, ioutils.GetBoolSuggests(), "")
	// cs.Properties.Enabled, err = strconv.ParseBool(enabled)
	if err != nil {
		return "", nil
	}
	answer += "," + ioutils.AskFromList("", "Insert the value for source.originAbsenceDetection >", false, ioutils.GetBoolSuggests(), "")
	// cs.Source.OriginAbsenceDetection, err = strconv.ParseBool(enabled)
	if err != nil {
		return "", nil
	}
	iq.AnswersMap[ContentSynchronisation] = answer
	return "", nil
}

// Specific writers for repo templates, since all the values in the templates should be written as string
var BoolToStringQuestionInfo = ioutils.QuestionInfo{
	Options:   ioutils.GetBoolSuggests(),
	AllowVars: true,
	Writer:    ioutils.WriteStringAnswer,
}

var IntToStringQuestionInfo = ioutils.QuestionInfo{
	Options:   nil,
	AllowVars: true,
	Writer:    ioutils.WriteStringAnswer,
}

var StringListToStringQuestionInfo = ioutils.QuestionInfo{
	Msg:       ioutils.CommaSeparatedListMsg,
	Options:   nil,
	AllowVars: true,
	Writer:    ioutils.WriteStringAnswer,
}

var questionMap = map[string]ioutils.QuestionInfo{
	TemplateType: {
		Options: []prompt.Suggest{
			{Text: Create, Description: "Template for creating a new repository"},
			{Text: Update, Description: "Template for updating an existing repository"},
		},
		Msg:          "",
		PromptPrefix: "Select the template type" + ioutils.PressTabMsg,
		AllowVars:    false,
		Writer:       ioutils.WriteStringAnswer,
		MapKey:       TemplateType,
		Callback:     nil,
	},
	ioutils.OptionalKey: {
		Msg:          "",
		PromptPrefix: SelectConfigKeyMsg,
		AllowVars:    false,
		Writer:       nil,
		MapKey:       "",
		Callback:     ioutils.OptionalKeyCallback,
	},
	Key: {
		Msg:          "",
		PromptPrefix: "Insert the repository key >",
		AllowVars:    true,
		Writer:       ioutils.WriteStringAnswer,
		MapKey:       Key,
		Callback:     nil,
	},
	Rclass: {
		Options: []prompt.Suggest{
			{Text: Local, Description: "A physical, locally-managed repository into which you can deploy artifacts"},
			{Text: Remote, Description: "A caching proxy for a repository managed at a remote URL"},
			{Text: Virtual, Description: "An Aggregation of several repositories with the same package type under a common URL."},
			{Text: Federated, Description: "A Federation is a collection of repositories of Federated type in different JPDs that are automatically configured for full bi-directional mirroring."},
		},
		Msg:          "",
		PromptPrefix: "Select the repository class" + ioutils.PressTabMsg,
		AllowVars:    false,
		Writer:       ioutils.WriteStringAnswer,
		MapKey:       Rclass,
		Callback:     rclassCallback,
	},
	MandatoryUrl: {
		Msg:          "",
		PromptPrefix: "Insert the remote repository URL >",
		AllowVars:    true,
		Writer:       ioutils.WriteStringAnswer,
		MapKey:       Url,
		Callback:     nil,
	},
	Url:             ioutils.FreeStringQuestionInfo,
	Description:     ioutils.FreeStringQuestionInfo,
	Notes:           ioutils.FreeStringQuestionInfo,
	IncludePatterns: StringListToStringQuestionInfo,
	ExcludePatterns: StringListToStringQuestionInfo,
	RepoLayoutRef: {
		Options: []prompt.Suggest{
			{Text: BowerDefaultRepoLayout},
			{Text: buildDefaultRepoLayout},
			{Text: ComposerDefaultRepoLayout},
			{Text: ConanDefaultRepoLayout},
			{Text: GoDefaultRepoLayout},
			{Text: GradleDefaultRepoLayout},
			{Text: IvyDefaultRepoLayout},
			{Text: Maven1DefaultRepoLayout},
			{Text: Maven2DefaultRepoLayout},
			{Text: NpmDefaultRepoLayout},
			{Text: NugetDefaultRepoLayout},
			{Text: puppetDefaultRepoLayout},
			{Text: SbtDefaultRepoLayout},
			{Text: SimpleDefaultRepoLayout},
			{Text: VcsDefaultRepoLayout},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	ProjectKey: {
		Options:   nil,
		AllowVars: false,
		Writer:    ioutils.WriteStringAnswer,
		Callback:  projectKeyCallback,
	},
	Environment: {
		PromptPrefix: "Insert the name of the environment to assign to >",
		AllowVars:    true,
		MapKey:       environmentsKey,
		Writer:       ioutils.WriteStringAnswer,
	},
	HandleReleases:               BoolToStringQuestionInfo,
	HandleSnapshots:              BoolToStringQuestionInfo,
	MaxUniqueSnapshots:           IntToStringQuestionInfo,
	SuppressPomConsistencyChecks: BoolToStringQuestionInfo,
	BlackedOut:                   BoolToStringQuestionInfo,
	DownloadRedirect:             BoolToStringQuestionInfo,
	BlockPushingSchema1:          BoolToStringQuestionInfo,
	DebianTrivialLayout:          BoolToStringQuestionInfo,
	ExternalDependenciesEnabled:  BoolToStringQuestionInfo,
	ExternalDependenciesPatterns: StringListToStringQuestionInfo,
	ChecksumPolicyType: {
		Options: []prompt.Suggest{
			{Text: ClientChecksumPolicy},
			{Text: ServerGeneratedChecksumsPolicy},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	MaxUniqueTags: IntToStringQuestionInfo,
	SnapshotVersionBehavior: {
		Options: []prompt.Suggest{
			{Text: UniqueBehavior},
			{Text: NonUniqueBehavior},
			{Text: DeployerBehavior},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	XrayIndex:              BoolToStringQuestionInfo,
	PropertySets:           StringListToStringQuestionInfo,
	ArchiveBrowsingEnabled: BoolToStringQuestionInfo,
	CalculateYumMetadata:   BoolToStringQuestionInfo,
	YumRootDepth:           IntToStringQuestionInfo,
	DockerApiVersion: {
		Options: []prompt.Suggest{
			{Text: DockerApiV1},
			{Text: DockerApiV2},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	EnableFileListsIndexing: BoolToStringQuestionInfo,
	OptionalIndexCompressionFormats: {
		Msg:       "Enter a comma separated list of values from " + strings.Join([]string{Bz2Compression, LzmaCompression, XzCompression}, ","),
		Options:   nil,
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	PrimaryKeyPairRef: ioutils.FreeStringQuestionInfo,
	Username: ioutils.FreeStringQuestionInfo,
	Password: ioutils.FreeStringQuestionInfo,
	Proxy:    ioutils.FreeStringQuestionInfo,
	RemoteRepoChecksumPolicyType: {
		Options: []prompt.Suggest{
			{Text: GenerateIfAbsentPolicy},
			{Text: FailPolicy},
			{Text: IgnoreAndGeneratePolicy},
			{Text: PassThruPolicy},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	HardFail:                          BoolToStringQuestionInfo,
	Offline:                           BoolToStringQuestionInfo,
	StoreArtifactsLocally:             BoolToStringQuestionInfo,
	SocketTimeoutMillis:               IntToStringQuestionInfo,
	LocalAddress:                      ioutils.FreeStringQuestionInfo,
	RetrievalCachePeriodSecs:          IntToStringQuestionInfo,
	FailedRetrievalCachePeriodSecs:    IntToStringQuestionInfo,
	MissedRetrievalCachePeriodSecs:    IntToStringQuestionInfo,
	UnusedArtifactsCleanupEnabled:     BoolToStringQuestionInfo,
	UnusedArtifactsCleanupPeriodHours: IntToStringQuestionInfo,
	AssumedOfflinePeriodSecs:          IntToStringQuestionInfo,
	FetchJarsEagerly:                  BoolToStringQuestionInfo,
	FetchSourcesEagerly:               BoolToStringQuestionInfo,
	RejectInvalidJars:                 BoolToStringQuestionInfo,
	ShareConfiguration:                BoolToStringQuestionInfo,
	SynchronizeProperties:             BoolToStringQuestionInfo,
	BlockMismatchingMimeTypes:         BoolToStringQuestionInfo,
	AllowAnyHostAuth:                  BoolToStringQuestionInfo,
	EnableCookieManagement:            BoolToStringQuestionInfo,
	BowerRegistryUrl:                  ioutils.FreeStringQuestionInfo,
	ComposerRegistryUrl:               ioutils.FreeStringQuestionInfo,
	PyPIRegistryUrl:                   ioutils.FreeStringQuestionInfo,
	VcsType: {
		Options: []prompt.Suggest{
			{Text: Git},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	VcsGitProvider: {
		Options: []prompt.Suggest{
			{Text: GithubVcsProvider},
			{Text: BitbucketVcsProvider},
			{Text: OldstashVcsProvider},
			{Text: StashVcsProvider},
			{Text: ArtifactoryVcsProvider},
			{Text: CustomVcsProvider},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	VcsGitDownloadUrl:         ioutils.FreeStringQuestionInfo,
	BypassHeadRequests:        BoolToStringQuestionInfo,
	ClientTlsCertificate:      ioutils.FreeStringQuestionInfo,
	FeedContextPath:           ioutils.FreeStringQuestionInfo,
	DownloadContextPath:       ioutils.FreeStringQuestionInfo,
	V3FeedUrl:                 ioutils.FreeStringQuestionInfo,
	ListRemoteFolderItems:     BoolToStringQuestionInfo,
	EnableTokenAuthentication: BoolToStringQuestionInfo,
	PodsSpecsRepoUrl:          ioutils.FreeStringQuestionInfo,
	ContentSynchronisation: {
		Options:   ioutils.GetBoolSuggests(),
		AllowVars: true,
		Writer:    nil,
		Callback:  contentSynchronisationCallBack,
	},
	Repositories: StringListToStringQuestionInfo,
	ArtifactoryRequestsCanRetrieveRemoteArtifacts: BoolToStringQuestionInfo,
	KeyPair: ioutils.FreeStringQuestionInfo,
	PomRepositoryReferencesCleanupPolicy: {
		Options: []prompt.Suggest{
			{Text: DiscardActiveRefrencePolicy},
			{Text: DiscardAnyReferencePolicy},
			{Text: NothingPolicy},
		},
		AllowVars: true,
		Writer:    ioutils.WriteStringAnswer,
	},
	DefaultDeploymentRepo:          ioutils.FreeStringQuestionInfo,
	ForceMavenAuthentication:       BoolToStringQuestionInfo,
	ForceNugetAuthentication:       BoolToStringQuestionInfo,
	ExternalDependenciesRemoteRepo: ioutils.FreeStringQuestionInfo,
}
