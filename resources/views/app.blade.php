<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="csrf-token" content="{{ csrf_token() }}">
        <meta name="author" content="Henry Whitaker">
        <meta name="version" content="{{ config('speedtest.version', 'Unknown') }}">

        <link href="/icons/themify/themify-icons.css" rel="stylesheet">

        <link rel="apple-touch-icon" sizes="57x57" href="/icons/fav/apple-icon-57x57.png">
        <link rel="apple-touch-icon" sizes="60x60" href="/icons/fav/apple-icon-60x60.png">
        <link rel="apple-touch-icon" sizes="72x72" href="/icons/fav/apple-icon-72x72.png">
        <link rel="apple-touch-icon" sizes="76x76" href="/icons/fav/apple-icon-76x76.png">
        <link rel="apple-touch-icon" sizes="114x114" href="/icons/fav/apple-icon-114x114.png">
        <link rel="apple-touch-icon" sizes="120x120" href="/icons/fav/apple-icon-120x120.png">
        <link rel="apple-touch-icon" sizes="144x144" href="/icons/fav/apple-icon-144x144.png">
        <link rel="apple-touch-icon" sizes="152x152" href="/icons/fav/apple-icon-152x152.png">
        <link rel="apple-touch-icon" sizes="180x180" href="/icons/fav/apple-icon-180x180.png">
        <link rel="icon" type="image/png" sizes="192x192"  href="/icons/fav/android-icon-192x192.png">
        <link rel="icon" type="image/png" sizes="32x32" href="/icons/fav/favicon-32x32.png">
        <link rel="icon" type="image/png" sizes="96x96" href="/icons/fav/favicon-96x96.png">
        <link rel="icon" type="image/png" sizes="16x16" href="/icons/fav/favicon-16x16.png">
        <link rel="manifest" href="/icons/fav/manifest.json">
        <meta name="msapplication-TileColor" content="#303030">
        <meta name="msapplication-TileImage" content="/icons/fav/ms-icon-144x144.png">
        <meta name="theme-color" content="#303030">

        <title>{{ config('app.name') }}</title>

        <!-- Styles -->
        <link rel="stylesheet" href="{{ mix('css/app.css') }}">

        <script src="{{ mix('js/app.js') }}"></script>
    </head>
    <body class="font-sans antialiased">
        @inertia
    </body>
</html>
