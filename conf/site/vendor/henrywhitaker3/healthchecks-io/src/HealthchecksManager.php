<?php

namespace Henrywhitaker3\Healthchecks;

use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksAccountLimitReachedException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksFailureException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksForbiddenException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUnauthorisedException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUuidNotFoundException;
use Henrywhitaker3\Healthchecks\Exceptions\InvalidUrlException;
use Henrywhitaker3\Healthchecks\Rules\IntValidationRule;
use Henrywhitaker3\Healthchecks\Rules\StringValidationRule;
use InvalidArgumentException;
use Rakit\Validation\Validator;
use Ramsey\Uuid\Exception\InvalidUuidStringException;
use Ramsey\Uuid\Uuid;

class HealthchecksManager
{
    /**
     * Healthchecks.io apikey
     *
     * @var String
     */
    private $apikey;

    /**
     * Healthchecks.io url
     *
     * @var String
     */
    private $url;

    /**
     * Authentication headers to send with requests
     *
     * @var array
     */
    private $headers;

    /**
     * Constructor for HealthchecksManager
     *
     * @param String $apikey
     * @param String $url
     */
    public function __construct(String $apikey, String $url = 'https://healthchecks.io/api/v1/')
    {
        $this->url = $url;
        $this->apikey = $apikey;
        $this->headers = [
            "X-Api-Key: $apikey",
        ];
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
	}

    /**
     * Get a list of existing checks
     *
     * @return array
     */
    public function listChecks()
    {
        $url = "$this->url/checks";

        $response = HttpClient::get($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return json_decode($response['body'], true)['checks'];
    }

    /**
     * Get a specific check
     *
     * @param String $uuid The UUID of the endpoint
     * @return array
     */
    public function getCheck(String $uuid)
    {
        if(Uuid::isValid($uuid) !== true) {
            throw new InvalidUuidStringException();
        }

        $url = "$this->url/checks/$uuid";

        $response = HttpClient::get($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            if($response['status'] == 403) {
                throw new HealthchecksForbiddenException($response['body']);
            }

            if($response['status'] == 404) {
                throw new HealthchecksUuidNotFoundException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return json_decode($response['body'], true);
    }

    /**
     * Pause a specific check
     *
     * @param String $uuid The UUID of the endpoint
     * @return bool
     */
    public function pauseCheck(String $uuid)
    {
        if(Uuid::isValid($uuid) !== true) {
            throw new InvalidUuidStringException();
        }

        $url = "$this->url/checks/$uuid/pause";

        $response = HttpClient::post($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            if($response['status'] == 403) {
                throw new HealthchecksForbiddenException($response['body']);
            }

            if($response['status'] == 404) {
                throw new HealthchecksUuidNotFoundException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return true;
    }

    /**
     * Resume a check (just pings the check)
     *
     * @param String $uuid The UUID of the endpoint
     * @return bool
     */
    public function resumeCheck(String $uuid)
    {
        if(Uuid::isValid($uuid) !== true) {
            throw new InvalidUuidStringException();
        }

        $url = "https://hc-ping.com/$uuid";

        $response = HttpClient::get($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            if($response['status'] == 403) {
                throw new HealthchecksForbiddenException($response['body']);
            }

            if($response['status'] == 404) {
                throw new HealthchecksUuidNotFoundException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return true;
    }

    /**
     * Delete a check
     *
     * @param String $uuid The UUID of the endpoint
     * @return bool
     */
    public function deleteCheck(String $uuid)
    {
        if(Uuid::isValid($uuid) !== true) {
            throw new InvalidUuidStringException();
        }

        $url = "$this->url/checks/$uuid";

        $response = HttpClient::delete($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            if($response['status'] == 403) {
                throw new HealthchecksForbiddenException($response['body']);
            }

            if($response['status'] == 404) {
                throw new HealthchecksUuidNotFoundException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return true;
    }

    /**
     * Get a specific check's pings
     *
     * @param String $uuid The UUID of the endpoint
     * @return array
     */
    public function getCheckPings(String $uuid)
    {
        if(Uuid::isValid($uuid) !== true) {
            throw new InvalidUuidStringException();
        }

        $url = "$this->url/checks/$uuid/pings";

        $response = HttpClient::get($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException();
            }

            if($response['status'] == 403) {
                throw new HealthchecksForbiddenException();
            }

            if($response['status'] == 404) {
                throw new HealthchecksUuidNotFoundException();
            }

            throw new HealthchecksFailureException();
        }

        return json_decode($response['body'], true);
    }

    /**
     * Get a specific check's status changes
     *
     * @param String $uuid The UUID of the endpoint
     * @return array
     */
    public function getCheckStatusChanges(String $uuid)
    {
        if(Uuid::isValid($uuid) !== true) {
            throw new InvalidUuidStringException();
        }

        $url = "$this->url/checks/$uuid/flips";

        $response = HttpClient::get($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            if($response['status'] == 403) {
                throw new HealthchecksForbiddenException($response['body']);
            }

            if($response['status'] == 404) {
                throw new HealthchecksUuidNotFoundException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return json_decode($response['body'], true);
    }

    /**
     * Create a new check
     *
     * @param array $args Array with optional parameters
     *   $params = [
     *      'name'         => (string) Name for check. Optional. 
     *                                   Default = ''.
     *      'tags'         => (string) Space-separated list of tags. Optional. 
     *                                   Default = ''.
     *      'desc'         => (string) Short description of check. Optional. 
     *                                   Default = null.
     *      'timeout'      => (int) Expected time between pings in seconds. Optional.
     *                                Default = 86400.
     *      'grace'        => (int) Time before alert after missed ping in seconds. Optional.
     *                                 Default = 3600.
     *      'schedule'     => (string) Cron expression for expected time b/w pings. Optional.
     *                                  Default = null.
     *      'tz'           => (string) Server timezone. Optional.
     *                                  Default = 'UTC'.
     *      'channels'      => (string) Comma-separated list of integration identifiers. Optional.
     *                                  Default = null.
     *    ]
     * @return array
     */
    public function createCheck(array $args = [])
    {
        if(!is_array($args)) {
            throw new InvalidArgumentException();
        }

        $rules = [
            'name' => 'string',
            'tags' => 'string',
            'desc' => 'string',
            'timeout' => 'int',
            'grace' => 'int',
            'schedule' => 'string',
            'tz' => 'string',
            'channels' => 'string',
        ];

        if(!$this->validator($args, $rules)) {
            throw new InvalidArgumentException();
        }

        $params = [
            'name' => '',
            'tags' => '',
            'desc' => null,
            'timeout' => 86400,
            'grace' => 3600,
            'schedule' => null,
            'tz' => 'UTC',
            'channels' => null,
        ];

        foreach($args as $key => $value) {
            $params[$key] = $value;
        }

        foreach($params as $key => $value) {
            if($value === null) {
                unset($params[$key]);
            }
        }

        $url = "$this->url/checks/";

        $response = HttpClient::post($url, $this->headers, $params);

        if($response['status'] !== 201 && $response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException();
            }

            if($response['status'] == 400) {
                throw new InvalidArgumentException($response['body']);
            }

            if($response['status'] == 403) {
                throw new HealthchecksAccountLimitReachedException();
            }
        }

        return json_decode($response['body'], true);
    }

    /**
     * Update a specific check
     *
     * @param String $uuid The UUID of the endpoint
     * @return array
     */
    public function updateCheck(String $uuid, array $args)
    {
        if(Uuid::isValid($uuid) !== true) {
            throw new InvalidUuidStringException();
        }

        if(!is_array($args)) {
            throw new InvalidArgumentException();
        }

        $rules = [
            'name' => 'string',
            'tags' => 'string',
            'desc' => 'string',
            'timeout' => 'int',
            'grace' => 'int',
            'schedule' => 'string',
            'tz' => 'string',
            'channels' => 'string',
        ];

        if(!$this->validator($args, $rules)) {
            throw new InvalidArgumentException();
        }

        $currentData = $this->getCheck($uuid);

        foreach($args as $key => $val) {
            $currentData[$key] = $val;
        }

        foreach($currentData as $key => $val) {
            if(!array_key_exists($key, $rules)) {
                unset($currentData[$key]);
            }
        }

        if(!isset($args['schedule']) && isset($args['timeout'])) {
            unset($currentData['schedule']);
        }

        $url = "$this->url/checks/$uuid";

        $response = HttpClient::post($url, $this->headers, $currentData);

        if($response['status'] !== 201 && $response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            if($response['status'] == 400) {
                throw new InvalidArgumentException($response['body']);
            }

            if($response['status'] == 403) {
                throw new HealthchecksAccountLimitReachedException($response['body']);
            }

            if($response['status'] == 404) {
                throw new HealthchecksUuidNotFoundException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return json_decode($response['body'], true);
    }

    /**
     * Get a list of intergations
     *
     * @return array
     */
    public function integrations()
    {

        $url = "$this->url/channels";

        $response = HttpClient::get($url, $this->headers);

        if($response['status'] !== 200) {
            if($response['status'] == 401) {
                throw new HealthchecksUnauthorisedException($response['body']);
            }

            throw new HealthchecksFailureException($response['body']);
        }

        return json_decode($response['body'], true);
    }

    /**
     * Validate an array
     *
     * @param array $array Input array
     * @param array $rules
     * @return bool
     */
    private function validator($array, $rules)
    {
        $validator = new Validator();
        $validator->addValidator('string', new StringValidationRule());
        $validator->addValidator('int', new IntValidationRule());
        $validation = $validator->validate($array, $rules);
        if($validation->fails()) {
            return false;
        }

        return true;
    }
}