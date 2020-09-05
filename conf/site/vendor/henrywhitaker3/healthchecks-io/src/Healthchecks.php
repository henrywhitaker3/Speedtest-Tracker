<?php

namespace Henrywhitaker3\Healthchecks;

use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksFailureException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUuidNotFoundException;
use Henrywhitaker3\Healthchecks\Exceptions\InvalidUrlException;
use Ramsey\Uuid\Exception\InvalidUuidStringException;
use Ramsey\Uuid\Uuid;

class Healthchecks
{
	/**
     * UUID for healthcheks
     *
     * @var string
     */
    public $uuid;

    /**
     * URL for healthcheks instance
     *
     * @var string
     */
    public $url;

    /**
     * Constructor for healthchecks class
     *
     * @param String $uuid
     * @param String $url
     */
    public function __construct(String $uuid, String $url = 'https://hc-ping.com/')
    {
        $this->uuid = $uuid;
        $this->url = $url;

        $this->validate();
	}

	/**
     * Validate the args supplied in the constructor
     *
     * @return void
     */
    public function validate()
    {
        if(substr($this->url, -1) != '/') {
            $this->url = $this->url . '/';
		}

		if(filter_var($this->url, FILTER_VALIDATE_URL) == false) {
			throw new InvalidUrlException();
		}

        if(Uuid::isValid($this->uuid) !== true) {
            throw new InvalidUuidStringException();
        }
    }

	/**
     * Send a 'success' signal
     *
     * @return boolean
     */
    public function success()
    {
		$url = $this->url . $this->uuid;

		$resp = HttpClient::get($url);

        if($resp['status'] !== 200) {
            if($resp['status'] == 404) {
                throw new HealthchecksUuidNotFoundException();
            }

            throw new HealthchecksFailureException();
        }

        return true;
	}

	/**
     * Send a 'fail' signal
     *
     * @return boolean
     */
    public function fail()
    {
		$url = $this->url . $this->uuid . '/fail';

		$resp = HttpClient::get($url);

        if($resp['status'] !== 200) {
            if($resp['status'] == 404) {
                throw new HealthchecksUuidNotFoundException();
            }

            throw new HealthchecksFailureException();
        }

        return true;
	}

	/**
     * Send a 'start' signal
     *
     * @return boolean
     */
    public function start()
    {
		$url = $this->url . $this->uuid . '/start';

		$resp = HttpClient::get($url);

        if($resp['status'] !== 200) {
            if($resp['status'] == 404) {
                throw new HealthchecksUuidNotFoundException();
            }

            throw new HealthchecksFailureException();
        }

        return true;
	}
}
