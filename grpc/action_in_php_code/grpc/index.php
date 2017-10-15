<?php
require dirname(__FILE__).'/vendor/autoload.php';
include_once 'GPBMetadata/Redis.php';
include_once 'Redis/RedisClient.php';
include_once 'Redis/RedisRequest.php';
include_once 'Redis/RedisReply.php';


$client = new Redis\RedisClient('localhost:50051', [
        'credentials' => Grpc\ChannelCredentials::createInsecure(),
]);

$request = new Redis\RedisRequest();
$request->setAction("get");
$request->setParam("key");
list($reply, $status) = $client->Command($request)->wait();
$message = $reply->getResult();
echo $message;
