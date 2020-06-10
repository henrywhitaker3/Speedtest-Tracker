<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="csrf-token" content="{{ csrf_token() }}">
        <meta name="author" content="Henry Whitaker">
        <meta name="version" content="{{ config('speedtest.version', 'Unknown') }}">

        <link href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/themify/themify-icons.css" rel="stylesheet">
        <link rel="stylesheet" href="{{ App\Helpers\SettingsHelper::getBase() }}files/css/bootstrap.dark.min.css">
        <link rel="stylesheet" href="{{ App\Helpers\SettingsHelper::getBase() }}files/css/main.css?v={{ str_replace('.', '-', config('speedtest.version')) }}">

        <link rel="apple-touch-icon" sizes="57x57" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-57x57.png">
        <link rel="apple-touch-icon" sizes="60x60" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-60x60.png">
        <link rel="apple-touch-icon" sizes="72x72" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-72x72.png">
        <link rel="apple-touch-icon" sizes="76x76" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-76x76.png">
        <link rel="apple-touch-icon" sizes="114x114" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-114x114.png">
        <link rel="apple-touch-icon" sizes="120x120" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-120x120.png">
        <link rel="apple-touch-icon" sizes="144x144" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-144x144.png">
        <link rel="apple-touch-icon" sizes="152x152" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-152x152.png">
        <link rel="apple-touch-icon" sizes="180x180" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/apple-icon-180x180.png">
        <link rel="icon" type="image/png" sizes="192x192"  href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/android-icon-192x192.png">
        <link rel="icon" type="image/png" sizes="32x32" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/favicon-32x32.png">
        <link rel="icon" type="image/png" sizes="96x96" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/favicon-96x96.png">
        <link rel="icon" type="image/png" sizes="16x16" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/favicon-16x16.png">
        <link rel="manifest" href="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/manifest.json">
        <meta name="msapplication-TileColor" content="#303030">
        <meta name="msapplication-TileImage" content="{{ App\Helpers\SettingsHelper::getBase() }}files/icons/fav/ms-icon-144x144.png">
        <meta name="theme-color" content="#303030">

        <title>{{ $title }}</title>
    </head>
    <body>
        <div id="main"></div>

        <script src="{{ App\Helpers\SettingsHelper::getBase() }}files/js/jquery.min.js"></script>
        <script src="{{ App\Helpers\SettingsHelper::getBase() }}files/js/popper.min.js"></script>
        <script src="{{ App\Helpers\SettingsHelper::getBase() }}files/js/app.js?v={{ str_replace('.', '-', config('speedtest.version')) }}"></script>
    </body>
</html>
