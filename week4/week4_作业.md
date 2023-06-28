**1.** 创建一个新的ibc数据包，名称为updatePost，包含postID、title、content。通过将该数据包发送给对手链，可更新存储在对手链中某条Post（通过postID）的标题和内容。
   对手链在ack确认数据包中将返回是否更新成功，发送链收到确认数据包后更新本地之前存储的SentPost为新的title。

```shell
ignite scaffold packet updatePost postID title content --ack postID title --module blog
```