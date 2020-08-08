<?php

namespace App\Helpers;

use Exception;
use File;
use Illuminate\Support\Facades\Log;
use RecursiveDirectoryIterator;
use RecursiveIteratorIterator;
use ZipArchive;

class UpdateHelper {
    /**
     * URL of updates
     *
     * @var string|bool
     */
    public $url;

    /**
     * Current app version number
     *
     * @var string|bool
     */
    public $currentVersion;

    /**
     * Latest app version number
     *
     * @var string|bool
     */
    public $latestVersion;

    /**
     * Username of GitHub repo
     *
     * @var string|bool
     */
    public $user;

    /**
     * Name of GitHub repo
     *
     * @var string|bool
     */
    public $repo;

    /**
     * Branch of GitHub repo
     *
     * @var string|bool
     */
    public $branch;

    /**
     * Store download
     *
     * @var string|null
     */
    public $download;

    function __construct() {
        $this->currentVersion = config('speedtest.version');
        $this->user = config('speedtest.user');
        $this->repo = config('speedtest.repo');
        $this->branch = config('speedtest.branch');
        $this->latestVersion = 'unknown';
        $this->download = null;
    }

    /**
     * Returns data on new version available
     *
     * @return  boolean|array   false|[ version, changelog ]
     */
    public function check()
    {
        Log::info('Checking for new version');

        if($this->currentVersion === false) {
            return false;
        }

        $gitVersion = $this->checkLatestVersion();
        if($gitVersion === false) {
            return false;
        }

        if((bool)(version_compare($this->currentVersion, $gitVersion['version']))) {
            Log::info('New version found! v' . $gitVersion['version']);
            $changelog = $this->getChangelog();
            return [
                'version' => $gitVersion['version'],
                'changelog' => $changelog[$gitVersion['version']],
            ];
        } else {
            return false;
        }
    }

    /**
     * Gets the latest version number from GitHub
     *
     * @return  array|bool   [ repo, branch, version ]
     */
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

    /**
     * Gets the latest changelog from GitHub
     *
     * @return  array
     */
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

    /**
     * Downloads the latest version from GitHub
     *
     * @return  boolean|Exception
     */
    public function downloadLatest()
    {
        Log::info('Downloading the latest version from GitHub');
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
            Log::info('New version successfully downloaded');
            return true;
        } catch(Exception $e) {
            Log::error('Couldn\'t download the update');
            Log::error($e);
            return $e;
        }
    }

    /**
     * Extracts zip archive from update
     *
     * @return  boolean
     */
    public function extractFiles()
    {
        Log::info('Extracting the update');
        $zip = new ZipArchive();
        $res = $zip->open('/tmp/'.$this->repo.'-update.zip');
        if($res === true) {
            $zip->extractTo('/tmp/'.$this->repo.'-update/');
            $zip->close();
            Log::info('Update extracted');
            return true;
        } else {
            Log::error('Couldn\'t extract the update');
            return false;
        }
    }

    /**
     * Replace existing files with newly downloaded files
     *
     * @return  void
     */
    public function updateFiles()
    {
        Log::info('Applying update');
        $dir = array_filter(glob('/tmp/'.$this->repo.'-update/*'), 'is_dir');
        $dir = $dir[0].DIRECTORY_SEPARATOR;

        $this->deleteExcluded($dir);
        $this->backupCurrent();
        $this->moveFiles();
        $this->clearup();

        Log::info('Successfully applied update');
    }

    /**
     * Deletes default templates from updated files.
     * This is for things like .env so that user specified files are not
     * overwritten.
     *
     * @param   string  $path
     * @return  void
     */
    private function deleteExcluded($path)
    {
        Log::info('Deleting excluded items from update directory');
        $exclude_dirs = config('speedtest.exclude_dirs', []);
        foreach($exclude_dirs as $dir) {
            $dir = $path . $dir;
            Log::debug('Deleting excluded directory: ' . $dir);

            File::deleteDirectory($dir);
        }

        $exclude_files = config('speedtest.exclude_files', []);
        foreach($exclude_files as $file) {
            $file = $path . $file;
            Log::debug('Deleting excluded file: ' . $file);

            File::delete($file);
        }
        Log::info('Excluded items deleted from update directory');
    }

    /**
     * Creates a ZIP backup of current installation
     *
     * @return  void
     */
    private function backupCurrent()
    {
        Log::info('Backing up current installation');

        $rootPath = realpath(base_path());
        $backupZip = '/tmp/speedtest-backup-'.time().'.zip';

        // Initialize archive object
        $zip = new ZipArchive();
        $zip->open($backupZip, ZipArchive::CREATE | ZipArchive::OVERWRITE);

        // Create recursive directory iterator
        $files = new RecursiveIteratorIterator(
            new RecursiveDirectoryIterator($rootPath),
            RecursiveIteratorIterator::LEAVES_ONLY
        );

        foreach ($files as $name => $file)
        {
            // Skip directories (they would be added automatically)
            if (!$file->isDir())
            {
                // Get real and relative path for current file
                $filePath = $file->getRealPath();
                $relativePath = substr($filePath, strlen($rootPath) + 1);

                // Add current file to archive
                $zip->addFile($filePath, $relativePath);
            }
        }

        // Zip archive will be created only after closing object
        $zip->close();
        Log::info('Backup created at: ' . $backupZip);
    }

    /**
     * Move updated files into server dir.
     *
     * @return  void
     */
    private function moveFiles()
    {
            $new = array_filter(glob('/tmp/'.$this->repo.'-update/*'), 'is_dir');
            $new = $new[0].DIRECTORY_SEPARATOR;

            foreach(File::files($new) as $file) {
                $filename = explode('/', $file);
                $filename = array_slice($filename, -1)[0];
                try {
                    Log::info('Overwriting ' . $filename);
                    Log::debug('From: ' . $file . ' to: ' . base_path().DIRECTORY_SEPARATOR.$filename);
                    File::delete(base_path().DIRECTORY_SEPARATOR.$filename);
                    File::move($file, base_path().DIRECTORY_SEPARATOR.$filename);
                } catch(Exception $e) {
                    Log::error('Failed to overwrite: ' . $filename);
                    Log::debug($e);
                }
            }

            $this->tempStoreExcludedFiles();

            foreach(File::directories($new) as $dir) {
                $dirname = explode('/', $dir);
                $dirname = array_slice($dirname, -1)[0];
                Log::info('Overwriting ' . $dir);
                File::deleteDirectory(base_path().DIRECTORY_SEPARATOR.$dirname);
                File::move($dir, base_path().DIRECTORY_SEPARATOR.$dirname);
            }

            $this->restoreExcludedFiles();

    }

    /**
     * Make a copy of excluded, user customised files
     *
     * @return  void
     */
    private function tempStoreExcludedFiles()
    {
        Log::info('Temporarily moving exluded files from root directory');
        foreach(config('speedtest.exclude_files', []) as $file) {
            try {
                Log::info('Moving ' . $file);
                File::copy(base_path().DIRECTORY_SEPARATOR.$file, '/tmp/'.$file);
            } catch(Exception $e) {
                Log::error('Couldn\'t backup '.$file);
            }
        }
        foreach(config('speedtest.exclude_dirs', []) as $dir) {
            try {
                Log::info('Moving ' . $dir);
                File::copyDirectory(base_path().DIRECTORY_SEPARATOR.$dir, '/tmp/'.$dir);
            } catch(Exception $e) {
                Log::error('Couldn\'t backup '.$dir);
            }
        }
    }

    /**
     * Restore user cusotmised files from the copy
     *
     * @return  void
     */
    private function restoreExcludedFiles()
    {
        Log::info('Restoring exluded files to root directory');
        foreach(config('speedtest.exclude_files', []) as $file) {
            try {
                Log::info('Moving ' . $file);
                File::copy('/tmp/'.$file, base_path().DIRECTORY_SEPARATOR.$file);
            } catch(Exception $e) {
                Log::error('Couldn\'t restore '.$file);
            }
        }
        foreach(config('speedtest.exclude_dirs', []) as $dir) {
            try {
                Log::info('Moving ' . $dir);
                File::copyDirectory('/tmp/'.$dir, base_path().DIRECTORY_SEPARATOR.$dir);
            } catch(Exception $e) {
                Log::error('Couldn\' restore ' . $dir);
            }
        }
    }

    /**
     * Delete update files from download dir.
     *
     * @return  void
     */
    private function clearup()
    {
        try {
            File::deleteDirectory('/tmp/'.$this->repo.'-update/');
            Log::info('Deleted download directory');
        } catch(Exception $e) {
            Log::error('Failed cleaning up update');
            Log::error($e);
        }
    }
}
