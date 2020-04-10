<?php

namespace App\Helpers;

use Exception;

class UpdateHelper {
    public static function check()
    {
        $current = config('speedtest.version', false);
        if($current === false) {
            return false;
        }

        $gitVersion = UpdateHelper::checkLatestVersion();
        if($gitVersion === false) {
            return false;
        }

        return (bool)(version_compare($current, $gitVersion['version']));
    }

    public static function checkLatestVersion()
    {
        $user = config('speedtest.user');
        $repo = config('speedtest.repo');
        $branch = config('speedtest.branch');

        $url = 'https://raw.githubusercontent.com/'.$user
                                                   .'/'
                                                   .$repo
                                                   .'/'
                                                   .$branch
                                                   .'/config/speedtest.php';

        try {
            $gitFile = file_get_contents($url);
        } catch(Exception $e) {
            return false;
        }

        $pattern = "/'version' => '([0-9]{1,}\.[0-9]{1,}\.[0-9]{1,})'/";
        $version = [];
        preg_match($pattern, $gitFile, $version);
        $version = $version[1];

        return [
            'repo' => $user . '/' . $repo,
            'branch' => $branch,
            'version' => $version,
        ];
    }
}
