<?php

use App\Helpers\SettingsHelper;
use GuzzleHttp\Psr7\MimeType;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Response;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

Route::get(SettingsHelper::getBase() . 'files/{path?}', function ($file) {
    $fileP = explode('?', $file)[0];
    $fileP = public_path() . '/' . $fileP;
    if (file_exists($fileP)) {
        $contents = File::get($fileP);
        $mime = MimeType::fromFilename($fileP);
        return Response::make(File::get($fileP), 200, ['Content-type' => $mime]);
    } else {
        abort(404);
    }
})->where('path', '.*')
    ->name('files');

Route::get('/{path?}', function () {
    return view('app', ['title' => SettingsHelper::get('app_name')->value]);
})->where('path', '^((?!\/api\/).)*$')
    ->name('react');
