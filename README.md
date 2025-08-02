# pipe2Tel
>I've only implemented the _sendMessage_ method, which fulfills my requirements for another project.
>The pipe2Tel tool currently supports only the _MarkdownV2_ parse_mode, but I'll enhance it as my needs evolve in different scenarios.

### Useage
```plaintext
Usage:
I>   pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID> [-restricted] [-msg=<TEXT OR FILE_PATH>]
II>  echo "sth" | pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID> [-restricted]

Options:
  -bot_token    The Telegram bot token (required)
  -chat_id      The Telegram chat ID (required)
  -rs           Optional flag to enable restricted mode (no web page preview, no notification)
  -msg          The message to send. If this is a file path, the file content is used as the message.
                If it's not a file path, it's treated as direct text.

If no -msg flag is provided, the program will read the message from stdin(II).
```

You may find these resources useful:
- https://core.telegram.org/bots/api
- [BotFather on t.me](https://t.me/botfather)
