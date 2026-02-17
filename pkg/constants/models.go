// Package constants provides shared constants across the codebase.
package constants

type ModelSpecs struct {
	ContextWindow   int
	MaxOutputTokens int
}

var ModelDefaults = map[string]ModelSpecs{
	"gpt-4o":                    {128000, 16384},
	"gpt-4o-mini":               {128000, 16384},
	"gpt-4-turbo":               {128000, 4096},
	"gpt-4-turbo-preview":       {128000, 4096},
	"gpt-4":                     {8192, 8192},
	"gpt-4-32k":                 {32768, 32768},
	"gpt-3.5-turbo":             {16384, 4096},
	"gpt-3.5-turbo-16k":         {16384, 4096},
	"o1":                        {200000, 100000},
	"o1-preview":                {128000, 32768},
	"o1-mini":                   {128000, 65536},
	"o3":                        {200000, 100000},
	"o3-mini":                   {200000, 100000},
	"o4-mini":                   {200000, 100000},
	"gpt-4.1":                   {1048576, 32768},
	"gpt-4.1-mini":              {1048576, 32768},

	"claude-opus-4":             {200000, 128000},
	"claude-opus-4-6":           {1000000, 128000},
	"claude-sonnet-4":           {200000, 16384},
	"claude-sonnet-4-5":         {1000000, 16384},
	"claude-haiku-4":            {200000, 8192},
	"claude-3-5-sonnet":         {200000, 8192},
	"claude-3-5-sonnet-20241022": {200000, 8192},
	"claude-3-5-haiku":          {200000, 8192},
	"claude-3-opus":             {200000, 4096},
	"claude-3-opus-20240229":    {200000, 4096},
	"claude-3-sonnet":           {200000, 4096},
	"claude-3-sonnet-20240229":  {200000, 4096},
	"claude-3-haiku":            {200000, 4096},
	"claude-3-haiku-20240307":   {200000, 4096},

	"gemini-2.5-pro":            {1048576, 65536},
	"gemini-2.5-flash":          {1048576, 65536},
	"gemini-2.5-flash-lite":     {1048576, 65536},
	"gemini-2.0-flash":          {1048576, 8192},
	"gemini-2.0-pro":            {2000000, 8192},
	"gemini-1.5-pro":            {2000000, 8192},
	"gemini-1.5-flash":          {1048576, 8192},

	"glm-4.7":                   {200000, 128000},
	"glm-4.7-flash":             {200000, 128000},
	"glm-4.7-flashx":            {200000, 128000},
	"glm-4.6":                   {200000, 128000},
	"glm-4":                     {128000, 4096},
	"glm-4-plus":                {128000, 4096},
	"glm-3-turbo":               {128000, 4096},

	"deepseek-chat":             {128000, 8000},
	"deepseek-reasoner":         {128000, 8000},
	"deepseek-coder":            {128000, 16384},

	"llama-3.1-8b-instant":      {131072, 131072},
	"llama-3.1-70b-versatile":   {131072, 131072},
	"llama-3.3-70b-versatile":   {131072, 32768},
	"llama-3.2-1b-preview":      {131072, 8192},
	"llama-3.2-3b-preview":      {131072, 8192},
	"llama-3.2-11b-vision-preview": {131072, 8192},
	"llama-3.2-90b-vision-preview": {131072, 8192},
	"mixtral-8x7b-32768":        {32768, 32768},
	"gemma2-9b-it":              {8192, 8192},
	"deepseek-r1-distill-llama-70b": {131072, 131072},

	"mistral-large-2407":        {128000, 128000},
	"mistral-large-3":           {128000, 128000},
	"mistral-medium-3":          {128000, 128000},
	"mistral-small-3.1":         {128000, 128000},
	"mistral-7b":                {128000, 8192},
	"codestral-latest":          {128000, 128000},
	"codestral-25-08":           {128000, 128000},
	"mixtral-8x7b":              {32768, 32768},
	"mixtral-8x22b":             {65536, 65536},
	"pixtral-large":             {128000, 128000},

	"meta-llama/Llama-3.1-8B":   {131072, 131072},
	"meta-llama/Llama-3.1-70B":  {131072, 131072},
	"meta-llama/Llama-3.1-405B": {131072, 131072},
	"meta-llama/Llama-3.2-1B":   {131072, 131072},
	"meta-llama/Llama-3.2-3B":   {131072, 131072},
	"meta-llama/Llama-3.2-11B-Vision": {131072, 131072},
	"meta-llama/Llama-3.2-90B-Vision": {131072, 131072},
	"meta-llama/Llama-3.3-70B":  {131072, 131072},
	"meta-llama/Llama-3-8B":     {8192, 2048},
	"meta-llama/Llama-3-70B":    {8192, 2048},

	"kimi-k1.5":                 {128000, 8192},
	"kimi-k2":                   {256000, 8192},
	"kimi-k2.5":                 {256000, 8192},
	"moonshot-v1-8k":            {8192, 4096},
	"moonshot-v1-32k":           {32768, 4096},
	"moonshot-v1-128k":          {131072, 4096},

	"minimax-m2.5":              {1000000, 1000000},
	"minimax-m2.5-lightning":    {1000000, 1000000},
	"minimax-m2":                {1000000, 1000000},
	"minimax-m2.1":              {1000000, 1000000},
	"minimax-2.5":               {24576, 8192},
}

func GetModelSpecs(model string) (ModelSpecs, bool) {
	specs, ok := ModelDefaults[model]
	return specs, ok
}

func GetModelContextWindow(model string) int {
	if specs, ok := ModelDefaults[model]; ok {
		return specs.ContextWindow
	}
	return 128000
}

func GetModelMaxOutputTokens(model string) int {
	if specs, ok := ModelDefaults[model]; ok {
		return specs.MaxOutputTokens
	}
	return 4096
}

func DefaultContextWindow() int   { return 128000 }
func DefaultMaxOutputTokens() int { return 4096 }
func DefaultTemperature() float64 { return 0.7 }
