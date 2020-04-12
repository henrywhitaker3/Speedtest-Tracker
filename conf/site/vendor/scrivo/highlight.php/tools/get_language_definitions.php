<?php

/* Copyright (c)
 * - 2013-2019  Geert Bergman (geert@scrivo.nl), highlight.php
 * - 2014       Daniel Lynge, highlight.php (contributor)
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

// Extract language definitions (JSON strings) from the large file that was
// created using 'node launcher.js' and create a JSON file for each language.

$f = file("languages.dat");

$patches = array(
    // The expression \\/: causes issues for PREG due to the / and the : having special meaning, therefore we use \Q and
    // \E to have PREG treat them as literal characters
    "1c" => array(
        array("\\\\\\\\/:", "\\\\Q\\\\/:\\\\E"),
    ),

    // The expression []{}%#'" should be treated as a list of invalid characters, however the [] is special in PREG so
    // we use \Q and \E to have PREG treat them as literal characters
    "ada" => array(
        array("[]{}%#'\\\"", "\\\\Q[]{}%#'\\\"\\\\E"),
    ),

    // WTF, any ideas anyone?
    "mercury" => array(array("\\\\\\\/", "\\\\\\\\\\\/")),

    // The expression [^] is not allowed in PREG
    "lisp" => array(array("[^]", "[^|]")),

    // There's a typo in the Swift translation file
    "swift" => array(array(
        '{02B80}-9',
        '{02B8}0-9',
    )),
);

for ($i = 0; $i < count($f); $i += 2) {
    if (!isset($f[$i + 1])) {
        continue;
    }

    $languageName = trim($f[$i]);
    $jsonLangDef = $f[$i + 1];

    // The `-` character must be escaped in while in `[]`. This is enforced in PHP 7.3+
    //   https://wiki.php.net/rfc/pcre2-migration
    //   https://github.com/php/php-src/pull/2857
    $jsonLangDef = preg_replace('/(\[[^:]*?\w)(-)([^a-zA-Z0-9\\\\]+?)/um', '$1\\\\\\-$3', $jsonLangDef);

    if (!$languageName) {
        die(sprintf("ERROR: No language name on line %d\n", ($i + 1)));
    }
    if (!@json_decode($jsonLangDef)) {
        die(sprintf("ERROR: Invalid JSON data on line %d\n", ($i + 2)));
    }

    if (isset($patches[$languageName])) {
        foreach ($patches[$languageName] as $j => $patch) {
            $patched = str_replace($patch[0], $patch[1], $jsonLangDef);

            if ($jsonLangDef === $patched) {
                printf("Patch %d for %s was not applied and likely unnecessary\n", $j, $languageName);
            }

            $jsonLangDef = $patched;
        }
    }

    $jsonLangDef = json_encode(json_decode($jsonLangDef), JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";

    if (!file_put_contents("../Highlight/languages/{$languageName}.json", $jsonLangDef)) {
        die("ERROR: Couldn't write to file.\n");
    }
}
