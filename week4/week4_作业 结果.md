**1.** 创建一个新的ibc数据包，名称为updatePost，包含postID、title、content。通过将该数据包发送给对手链，可更新存储在对手链中某条Post（通过postID）的标题和内容。
   对手链在ack确认数据包中将返回是否更新成功，发送链收到确认数据包后更新本地之前存储的SentPost为新的title。

```shell
ignite scaffold packet updatePost postID title content --ack postID title --module blog
```




**2.** 修改keeper方法中的`OnRecvUpdatePostPacket `。

```shell
	id, err := strconv.ParseUint(data.GetPostID(), 10, 64)
	if err != nil {
		return packetAck, err
	}
	post, found := k.GetPost(ctx, id)
	if !found {
		return packetAck, errors.New("post not found")
	}
	post.Title = data.GetTitle()
	post.Content = data.GetContent()
	k.SetPost(ctx, post)
	packetAck.PostID = data.PostID

```

**3.** 修改keeper方法中的`OnAcknowledgementUpdatePostPacket `。

```shell
     id, _ := strconv.ParseUint(data.PostID, 10, 64) // 字符串转换为 int
     post, found := k.GetSentPost(ctx, id)
     if !found {
         return errors.New("cannot find post")
     }
     post.Title = data.GetTitle()
     k.SetSentPost(ctx, post)
```



**4.** 分别启动两条链

```
ignite chain serve -c earth.yml

ignite chain serve -c mars.yml
```


**5.** 启动中继器

```
rm -rf ~/.ignite/relayer

ignite relayer configure -a \
  --source-rpc "http://0.0.0.0:26657" \
  --source-faucet "http://0.0.0.0:4500" \
  --source-port "blog" \
  --source-version "blog-1" \
  --source-gasprice "0.0000025stake" \
  --source-prefix "cosmos" \
  --source-gaslimit 300000 \
  --target-rpc "http://0.0.0.0:26659" \
  --target-faucet "http://0.0.0.0:4501" \
  --target-port "blog" \
  --target-version "blog-1" \
  --target-gasprice "0.0000025stake" \
  --target-prefix "cosmos" \
  --target-gaslimit 300000

ignite relayer connect
```


 