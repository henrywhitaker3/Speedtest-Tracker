<?php

/* Copyright (c) 2019 Geert Bergman (geert@scrivo.nl), highlight.php
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
 * A PHP representation of a Mode in the JS library.
 *
 * @internal
 *
 * @since 9.16.0.0
 * @mixin ModeDeprecations
 *
 * Language definition set via language definition JSON files
 *
 * @property bool $case_insensitive = false
 * @property string[] $aliases = array()
 * @property string|null $className = null
 * @property string|null $begin = null
 * @property RegEx|null $beginRe = null
 * @property string|null $end = null
 * @property RegEx|null $endRe = null
 * @property string|null $beginKeywords = null
 * @property bool $endsWithParent = false
 * @property bool $endsParent = false
 * @property bool $endSameAsBegin = false
 * @property string|null $lexemes = null
 * @property RegEx|null $lexemesRe = null
 * @property array<string, array<int, string|int>> $keywords = array()
 * @property string|null $illegal = null
 * @property RegEx|null $illegalRe = null
 * @property bool $excludeBegin = false
 * @property bool $excludeEnd = false
 * @property bool $returnBegin = false
 * @property bool $returnEnd = false
 * @property Mode[] $contains = array()
 * @property Mode|null $starts = null
 * @property Mode[] $variants = array()
 * @property int|null $relevance = null
 * @property string|string[]|null $subLanguage = null
 * @property bool $skip = false
 * @property bool $disableAutodetect = false
 *
 * Properties set at runtime by the language compilation process
 * @property array<int, Mode> $cachedVariants = array()
 * @property Terminators|null $terminators = null
 * @property string $terminator_end = ""
 * @property bool $compiled = false
 * @property Mode|null $parent = null
 * @property string $type = ''
 *
 * @see https://highlightjs.readthedocs.io/en/latest/reference.html
 */
abstract class Mode extends \stdClass
{
    /**
     * Fill in the missing properties that this Mode does not have.
     *
     * @internal
     *
     * @param \stdClass|null $obj
     *
     * @since 9.16.0.0
     *
     * @return void
     */
    public static function _normalize(&$obj)
    {
        // Don't overload our Modes if we've already normalized it
        if (isset($obj->__IS_COMPLETE)) {
            return;
        }

        if ($obj === null) {
            $obj = new \stdClass();
        }

        $patch = array(
            "begin" => true,
            "end" => true,
            "lexemes" => true,
            "illegal" => true,
        );

        // These values come in from JSON definition files
        $defaultValues = array(
            "case_insensitive" => false,
            "aliases" => array(),
            "className" => null,
            "begin" => null,
            "beginRe" => null,
            "end" => null,
            "endRe" => null,
            "beginKeywords" => null,
            "endsWithParent" => false,
            "endsParent" => false,
            "endSameAsBegin" => false,
            "lexemes" => null,
            "lexemesRe" => null,
            "keywords" => array(),
            "illegal" => null,
            "illegalRe" => null,
            "excludeBegin" => false,
            "excludeEnd" => false,
            "returnBegin" => false,
            "returnEnd" => false,
            "contains" => array(),
            "starts" => null,
            "variants" => array(),
            "relevance" => null,
            "subLanguage" => null,
            "skip" => false,
            "disableAutodetect" => false,
        );

        // These values are set at runtime
        $runTimeValues = array(
            "cachedVariants" => array(),
            "terminators" => null,
            "terminator_end" => "",
            "compiled" => false,
            "parent" => null,

            // This value is unique to highlight.php Modes
            "__IS_COMPLETE" => true,
        );

        foreach ($patch as $k => $v) {
            if (isset($obj->{$k})) {
                $obj->{$k} = str_replace("\\/", "/", $obj->{$k});
                $obj->{$k} = str_replace("/", "\\/", $obj->{$k});
            }
        }

        foreach ($defaultValues as $k => $v) {
            if (!isset($obj->{$k}) && is_object($obj)) {
                $obj->{$k} = $v;
            }
        }

        foreach ($runTimeValues as $k => $v) {
            if (is_object($obj)) {
                $obj->{$k} = $v;
            }
        }
    }

    /**
     * Set any deprecated properties values to their replacement values.
     *
     * @internal
     *
     * @param \stdClass $obj
     *
     * @return void
     */
    public static function _handleDeprecations(&$obj)
    {
        $deprecations = array(
            // @TODO Deprecated since 9.16.0.0; remove at 10.x
            'caseInsensitive' => 'case_insensitive',
            'terminatorEnd' => 'terminator_end',
        );

        foreach ($deprecations as $deprecated => $new) {
            $obj->{$deprecated} = &$obj->{$new};
        }
    }
}
