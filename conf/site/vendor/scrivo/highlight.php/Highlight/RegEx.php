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
 * A PHP implementation to match JavaScript's RegExp class as closely as possible.
 *
 * A lot of behavior in this class is reversed engineered, so improvements are welcome!
 *
 * @internal
 *
 * @since 9.16.0
 */
final class RegEx
{
    /**
     * @var string
     */
    public $source;

    /**
     * @var int
     */
    public $lastIndex = 0;

    /**
     * @param RegEx|string $regex
     */
    public function __construct($regex)
    {
        $this->source = (string) $regex;
    }

    public function __toString()
    {
        return (string) $this->source;
    }

    /**
     * Run the regular expression against the given string.
     *
     * @since 9.16.0.0
     *
     * @param string $str the string to run this regular expression against
     *
     * @return RegExMatch|null
     */
    public function exec($str)
    {
        $index = null;
        $results = array();
        preg_match_all($this->source, $str, $results, PREG_SET_ORDER | PREG_OFFSET_CAPTURE, $this->lastIndex);

        if ($results === null || count($results) === 0) {
            return null;
        }

        foreach ($results[0] as &$result) {
            if ($result[1] !== -1) {
                // Only save the index if it hasn't been set yet
                if ($index === null) {
                    $index = $result[1];
                }

                $result = $result[0];
            } else {
                $result = null;
            }
        }

        $results = $results[0];
        $this->lastIndex += mb_strlen($results[0]) + ($index - $this->lastIndex);

        $matches = new RegExMatch($results);
        $matches->index = isset($index) ? $index : 0;
        $matches->input = $str;

        return $matches;
    }
}
