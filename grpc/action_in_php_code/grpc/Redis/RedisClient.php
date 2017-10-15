<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Redis;

/**
 */
class RedisClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Redis\RedisRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function Command(\Redis\RedisRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/redis.Redis/Command',
        $argument,
        ['\Redis\RedisReply', 'decode'],
        $metadata, $options);
    }

}
