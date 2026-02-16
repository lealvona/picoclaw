# PicoClaw Provider Context Window Research

## First-Class Providers

### OpenAI-Compatible Providers
| Provider | Context Window | Max Output | Notes |
|----------|---------------|-------------|--------|
| OpenAI | 128K-128K | 4,096-16,384 | GPT-4o, GPT-4-turbo, etc. |
| Anthropic | 200K | 4,096 | Claude 3.5 Sonnet, Opus, Haiku |
| Groq | 128K | 8,192 | Llama 3.1 70B, Mixtral 8x7B |
| Gemini | 2M | 8,192 | Gemini 1.5 Pro, Flash |
| DeepSeek | 64K-128K | 2K-8K | DeepSeek-V3, Coder V2 |
| Zhipu (z.ai) | 128K | 4,096 | GLM-4, GLM-4-Air, GLM-4-Flash |

### Chinese Providers (no USA endpoints)
| Provider | Context Window | Max Output | Notes |
|----------|---------------|-------------|--------|
| Moonshot (z.ai) | 32K | 8,192 | Kimi models |
| Minimax | 24K-16K | 8,192 | MiniMax series |

### Self-Hosted / Custom
| Provider | Context Window | Max Output | Notes |
|----------|---------------|-------------|--------|
| VLLM | Variable | Variable | User-hosted models |
| Ollama | Variable | Variable | Local LLM inference |
| Nvidia | 128K | 4,096 | AI Foundation models |

### Recommended Settings

#### Minimax Configuration
- Context Window: 24,576 tokens (supports MiniMax 2.5)
- Max Output: 8,192 tokens
- Max Tokens: 24,576
- Temperature: 0.7 (default)
- API Base: https://api.minimax.chat/v1

#### Zhipu Configuration
- Context Window: 128,000 tokens
- Max Output: 4,096 tokens
- Max Tokens: 128,000
- Temperature: 0.7 (default)
- API Base: https://open.bigmodel.cn/api/paas/v4

## Models Available

### Minimax Models
- `minimax-2.5` - Main model, 24K context, high quality
- `minimax-2.5-128k` - Extended context, 128K tokens
- `mini-4` - Faster, smaller context
- `mini-5` - Standard model

### Zhipu (z.ai) Models
- `glm-4` - Latest flagship, 128K context
- `glm-4-air` - Optimized for coding tasks
- `glm-4-flash` - Faster, slightly lower quality
- `glm-4-plus` - Enhanced version
- `glm-3-turbo` - Faster, lower cost
- `glm-3` - Previous generation

### OpenAI Models (Lealvona Server)
The Lealvona server is OpenAI-compatible and exposes:
- OpenAI models via API
- Same interface as OpenAI
- Can be used with OpenAI provider configuration

## Configuration Notes

1. **Max Tokens**: Should match or be less than context window
2. **Max Output Tokens**: Typically 1/2 to 1/4 of context window
3. **Temperature**: 0.7 for balanced responses, 0.0-1.0 range
4. **Custom Provider**: For self-hosted or non-standard providers

## API Endpoints

### USA-Facing Endpoints
- OpenAI: https://api.openai.com/v1
- Anthropic: https://api.anthropic.com/v1
- Groq: https://api.groq.com/openai/v1
- Gemini: https://generativelanguage.googleapis.com
- DeepSeek: https://api.deepseek.com
- OpenRouter: https://openrouter.ai/api/v1
- Nvidia: https://integrate.api.nvidia.com/v1
- Minimax: https://api.minimax.chat/v1

### Chinese Providers (no USA endpoints)
- Zhipu (z.ai): https://open.bigmodel.cn/api/paas/v4
- Moonshot: https://api.moonshot.cn/v1

### Self-Hosted
- VLLM: User-defined
- Ollama: http://localhost:11434/v1 (default)

## Best Practices

1. **Context Window Matching**: Ensure max_tokens ≤ context_window
2. **Output Optimization**: max_output_tokens ≈ 1/3 to 1/2 of context
3. **Temperature Settings**: 
   - Creative tasks: 0.8-1.0
   - Coding tasks: 0.2-0.4
   - Analytical: 0.0-0.2
4. **Token Efficiency**: Use larger models for complex tasks, smaller for simple

## References

- OpenAI API Docs: https://platform.openai.com/docs
- Anthropic API Docs: https://docs.anthropic.com
- Zhipu (z.ai) API: https://open.bigmodel.cn/dev/api
- Minimax API: https://www.minimax.chat/document
