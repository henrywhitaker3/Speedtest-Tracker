<?php

namespace App\Http\Controllers;

use App\Helpers\SettingsHelper;
use Exception;
use Updater;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;

class UpdateController extends Controller
{
    public function __construct()
    {
        if((bool)SettingsHelper::get('auth')->value === true) {
            $this->middleware('auth:api');
        }
    }

    /**
     * Check for new update
     *
     * @return  JsonResponse
     */
    public function checkForUpdate()
    {
        return response()->json([
            'method' => 'check for updates',
            'update' => Updater::check(),
        ], 200);
    }

    /**
     * Download new update
     *
     * @return  JsonResponse
     */
    public function downloadUpdate()
    {
        $dl = Updater::downloadLatest();

        if($dl) {
            return response()->json([
                'method' => 'download latest version',
                'success' => true,
            ], 200);
        } else {
            return response()->json([
                'method' => 'download latest version',
                'success' => false,
            ], 500);
        }
    }

    /**
     * Trigger update extraction
     *
     * @return  JsonResponse
     */
    public function extractUpdate()
    {
        $ex = Updater::extractFiles();

        if($ex) {
            return response()->json([
                'method' => 'extract latest version',
                'success' => true,
            ], 200);
        } else {
            return response()->json([
                'method' => 'extract latest version',
                'success' => false,
            ], 500);
        }
    }

    /**
     * Trigger update file move
     *
     * @return  JsonResponse
     */
    public function moveUpdate()
    {
        Updater::updateFiles();

        return response()->json([
            'method' => 'copy latest version',
            'success' => null,
        ], 200);
    }

    /**
     * Get local changelog
     *
     * @return  JsonResponse
     */
    public function changelog()
    {
        $url = base_path() . '/changelog.json';

        try {
            $changelog = json_decode(file_get_contents($url), true);
        } catch(Exception $e) {
            $changelog = [];
        }

        return response()->json([
            'method' => 'get changelog',
            'data' => $changelog
        ], 200);
    }
}
