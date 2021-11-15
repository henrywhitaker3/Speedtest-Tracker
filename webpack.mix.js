const mix = require("laravel-mix");

require("laravel-mix-tailwind");


/*
 |--------------------------------------------------------------------------
 | Mix Asset Management
 |--------------------------------------------------------------------------
 |
 | Mix provides a clean, fluent API for defining some Webpack build steps
 | for your Laravel application. By default, we are compiling the Sass
 | file for the application as well as bundling up all the JS files.
 |
 */

mix.js("resources/js/app.js", "public/js")
    .vue()
    .webpackConfig(require("./webpack.config"))
    .version();

mix.sass("resources/sass/app.scss", "public/css")
    .options({
        processCssUrls: false,
        imgLoaderOptions: {
            enabled: false,
        },
    })
    .sass("resources/sass/tailwind-utilities.scss", "public/css/app.css")
    .tailwind("./tailwind.config.js")
    .version();
