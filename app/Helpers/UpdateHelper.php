<?php

namespace App\Helpers;

use Exception;
use ZipArchive;

class UpdateHelper {
    public $url;
    public $currentVersion;
    public $user;
    public $repo;
    public $branch;

    function __construct() {
        $this->currentVersion = config('speedtest.version');
        $this->user = config('speedtest.user');
        $this->repo = config('speedtest.repo');
        $this->branch = config('speedtest.branch');
        $this->latestVersion = 'unknown';
        $this->download = null;
    }

    public function check()
    {
        if($this->currentVersion === false) {
            return false;
        }

        $gitVersion = $this->checkLatestVersion();
        if($gitVersion === false) {
            return false;
        }

        if((bool)(version_compare($this->currentVersion, $gitVersion['version']))) {
            $changelog = $this->getChangelog();
            return [
                'version' => $gitVersion['version'],
                'changelog' => $changelog[$gitVersion['version']],
            ];
        } else {
            return false;
        }
    }

    public function checkLatestVersion()
    {

        $url = 'https://raw.githubusercontent.com/'
               .$this->user
               .'/'
               .$this->repo
               .'/'
               .$this->branch
               .'/config/speedtest.php';

        try {
            $gitFile = file_get_contents($url);
        } catch(Exception $e) {
            return false;
        }

        $pattern = "/'version' => '([0-9]{1,}\.[0-9]{1,}\.[0-9]{1,})'/";
        $version = [];
        preg_match($pattern, $gitFile, $version);
        $this->latestVersion = $version[1];

        return [
            'repo' => $this->user . '/' . $this->repo,
            'branch' => $this->branch,
            'version' => $this->latestVersion,
        ];
    }

    public function getChangelog()
    {
        $url = 'https://raw.githubusercontent.com/'
               .$this->user
               .'/'
               .$this->repo
               .'/'
               .$this->branch
               .'/changelog.json';

        try {
            $changelog = json_decode(file_get_contents($url), true);
        } catch(Exception $e) {
            $changelog = [];
        }

        return $changelog;
    }

    public function downloadLatest()
    {
        $url = 'https://github.com/'
                .$this->user
                .'/'
                .$this->repo
                .'/archive/'
                .$this->branch
                .'.zip';

        try {
            $zip = file_get_contents($url);
            $name = '/tmp/'.$this->repo.'-update.zip';
            file_put_contents($name, $zip);
            return true;
        } catch(Exception $e) {
            return $e;
        }
    }

    public function extractFiles()
    {
        $zip = new ZipArchive();
        $res = $zip->open('/tmp/'.$this->repo.'-update.zip');
        if($res === true) {
            $zip->extractTo('/tmp/'.$this->repo.'-update/');
            $zip->close();
            return true;
        } else {
            return false;
        }
    }

    public function updateFiles()
    {
        foreach (glob('/tmp/'.$this->repo.'-update/') as $folder) {
            return $folder;
        }
    }
}
