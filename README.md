# pipe2Tel

A command-line tool to send text or files into Telegram chats via a Bot API, with automatic handling of long messages.

## Installation

```bash
go install github.com/awirshf45d/pipe2Tel/cmd/pipe2Tel@vlatest
```

## Usage

```plaintext
I>   pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID> [-rs] -msg=<TEXT OR FILE_PATH>
II>  echo "sth" | pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID> [-rs]
```
* **`-bot_token`** (`string`, required): Your Telegram bot token. 
* **`-chat_id`** (`string`, required): Target chat ID or channel username.
* **`-msg`** (`string`, optional): Message text or file path. If omitted, reads from stdin.
* **`-rs`** (`flag`, optional): Restricted mode (no previews, no web page).

## Details

* **MarkdownV2** by default with full escaping of special characters (excluding backticks, underscore, asterisk) so code fences and bold/italic texts render correctly. It covers: backslash (`\`), square brackets (`[ ]`), parentheses (`( )`), tilde (`~`), greater-than (`>`), hash (`#`), plus (`+`), minus (`-`), equal (`=`), pipe (`|`), curly braces (`{ }`), period (`.`), and exclamation mark (`!`).
* Short messages use `sendMessage`; long texts & files use `sendDocument` with `multipart/form-data`.
* Temporary files are cleaned up automatically.

## Resources
You may find these resources useful:
- https://core.telegram.org/bots/api
- [BotFather on t.me](https://t.me/botfather)
