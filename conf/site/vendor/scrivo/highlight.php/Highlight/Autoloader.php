<?php

/* Copyright (c) 2013-2019 Geert Bergman (geert@scrivo.nl), highlight.php
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 * 3. Neither the name of "highlight.js", "highlight.php", nor the names of its
 *    contributors may be used to endorse or promote products derived from this
 *    software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

namespace Highlight;

/**
 * The autoloader class for Highlight classes for projects without Composer or autoloaders.
 *
 * ```php
 * <?php
 *
 * require_once("Highlight/Autoloader.php");
 * spl_autoload_register("\\Highlight\\Autoloader::load");
 *
 * // Now use Highlight classes:
 * $hl = new \Highlight\Highlighter(...);
 * ```
 *
 * @TODO In v10.x, mark this class final for real
 *
 * @api
 * @final
 *
 * @since 7.5.0.0
 */
class Autoloader
{
    /**
     * The method to include the source file for a given class to use in
     * the PHP spl_autoload_register function.
     *
     * @param string $class a name of a Scrivo class
     *
     * @return bool true if the source file was successfully included
     */
    public static function load($class)
    {
        if (substr($class, 0, 10) !== "Highlight\\") {
            return false;
        }

        $c = str_replace("\\", "/", substr($class, 10)) . ".php";
        $res = include __DIR__ . "/$c";

        return $res == 1 ? true : false;
    }
}
