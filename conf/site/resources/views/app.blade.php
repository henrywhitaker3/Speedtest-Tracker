<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="csrf-token" content="{{ csrf_token() }}">
        <meta name="author" content="Henry Whitaker">
        <meta name="version" content="{{ config('speedtest.version', 'Unknown') }}">

        <link href="/icons/themify/themify-icons.css" rel="stylesheet">
        <link rel="stylesheet" href="/css/bootstrap.dark.min.css">
        <link rel="stylesheet" href="/css/main.css">

        <title>{{ $title }}</title>
    </head>
    <body>
        <div id="main"></div>

        <script src="/js/jquery.min.js"></script>
        <script src="/js/popper.min.js"></script>
        <script src="/js/app.js"></script>
    </body>
</html>
