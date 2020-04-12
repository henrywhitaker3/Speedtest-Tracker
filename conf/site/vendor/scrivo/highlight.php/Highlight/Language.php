<?php

/* Copyright (c)
 * - 2006-2013, Ivan Sagalaev (maniacsoftwaremaniacs.org), highlight.js
 *              (original author)
 * - 2013-2019, Geert Bergman (geertscrivo.nl), highlight.php
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

namespace Highlight;

/**
 * @todo In highlight.php 10.x, replace the @final attribute with the `final` keyword.
 *
 * @final
 *
 * @internal
 *
 * // Backward compatibility properties
 *
 * @property Mode $mode (DEPRECATED) All properties traditionally inside of $mode are now available directly from this class.
 * @property bool $caseInsensitive (DEPRECATED) Due to compatibility requirements with highlight.js, use `case_insensitive` instead.
 */
class Language extends Mode
{
    /** @var string[] */
    private static $COMMON_KEYWORDS = array('of', 'and', 'for', 'in', 'not', 'or', 'if', 'then');

    /** @var string */
    public $name;

    /** @var Mode|null */
    private $mode = null;

    /**
     * @param string $lang
     * @param string $filePath
     *
     * @throws \InvalidArgumentException when the given $filePath is inaccessible
     */
    public function __construct($lang, $filePath)
    {
        $this->name = $lang;

        // We're loading the JSON definition file as an \stdClass object instead of an associative array. This is being
        // done to take advantage of objects being pass by reference automatically in PHP whereas arrays are pass by
        // value.
        $json = file_get_contents($filePath);

        if ($json === false) {
            throw new \InvalidArgumentException("Language file inaccessible: $filePath");
        }

        $this->mode = json_decode($json);
    }

    /**
     * @param string $name
     *
     * @return bool|Mode|null
     */
    public function __get($name)
    {
        if ($name === 'mode') {
            @trigger_error('The "mode" property will be removed in highlight.php 10.x', E_USER_DEPRECATED);

            return $this->mode;
        }

        if ($name === 'caseInsensitive') {
            @trigger_error('Due to compatibility requirements with highlight.js, use "case_insensitive" instead.', E_USER_DEPRECATED);

            if (isset($this->mode->case_insensitive)) {
                return $this->mode->case_insensitive;
            }

            return false;
        }

        if (isset($this->mode->{$name})) {
            return $this->mode->{$name};
        }

        return null;
    }

    /**
     * @param string $value
     * @param bool   $global
     *
     * @return RegEx
     */
    private function langRe($value, $global = false)
    {
        return RegExUtils::langRe($value, $global, $this->case_insensitive);
    }

    /**
     * Performs a shallow merge of multiple objects into one.
     *
     * @param Mode                 $params the objects to merge
     * @param array<string, mixed> ...$_
     *
     * @return Mode
     */
    private function inherit($params, $_ = array())
    {
        /** @var Mode $result */
        $result = new \stdClass();
        $objects = func_get_args();
        $parent = array_shift($objects);

        foreach ($parent as $key => $value) {
            $result->{$key} = $value;
        }

        foreach ($objects as $object) {
            foreach ($object as $key => $value) {
                $result->{$key} = $value;
            }
        }

        return $result;
    }

    /**
     * @param Mode|null $mode
     *
     * @return bool
     */
    private function dependencyOnParent($mode)
    {
        if (!$mode) {
            return false;
        }

        if (isset($mode->endsWithParent) && $mode->endsWithParent) {
            return $mode->endsWithParent;
        }

        return $this->dependencyOnParent(isset($mode->starts) ? $mode->starts : null);
    }

    /**
     * @param Mode $mode
     *
     * @return array<int, \stdClass|Mode>
     */
    private function expandOrCloneMode($mode)
    {
        if ($mode->variants && !$mode->cachedVariants) {
            $mode->cachedVariants = array();

            foreach ($mode->variants as $variant) {
                $mode->cachedVariants[] = $this->inherit($mode, array('variants' => null), $variant);
            }
        }

        // EXPAND
        // if we have variants then essentially "replace" the mode with the variants
        // this happens in compileMode, where this function is called from
        if ($mode->cachedVariants) {
            return $mode->cachedVariants;
        }

        // CLONE
        // if we have dependencies on parents then we need a unique
        // instance of ourselves, so we can be reused with many
        // different parents without issue
        if ($this->dependencyOnParent($mode)) {
            return array($this->inherit($mode, array(
                'starts' => $mode->starts ? $this->inherit($mode->starts) : null,
            )));
        }

        // highlight.php does not have a concept freezing our Modes

        // no special dependency issues, just return ourselves
        return array($mode);
    }

    /**
     * @param Mode      $mode
     * @param Mode|null $parent
     *
     * @return void
     */
    private function compileMode($mode, $parent = null)
    {
        Mode::_normalize($mode);

        if ($mode->compiled) {
            return;
        }

        $mode->compiled = true;
        $mode->keywords = $mode->keywords ? $mode->keywords : $mode->beginKeywords;

        if ($mode->keywords) {
            $mode->keywords = $this->compileKeywords($mode->keywords, (bool) $this->case_insensitive);
        }

        $mode->lexemesRe = $this->langRe($mode->lexemes ? $mode->lexemes : "\w+", true);

        if ($parent) {
            if ($mode->beginKeywords) {
                $mode->begin = "\\b(" . implode("|", explode(" ", $mode->beginKeywords)) . ")\\b";
            }

            if (!$mode->begin) {
                $mode->begin = "\B|\b";
            }

            $mode->beginRe = $this->langRe($mode->begin);

            if ($mode->endSameAsBegin) {
                $mode->end = $mode->begin;
            }

            if (!$mode->end && !$mode->endsWithParent) {
                $mode->end = "\B|\b";
            }

            if ($mode->end) {
                $mode->endRe = $this->langRe($mode->end);
            }

            $mode->terminator_end = $mode->end;

            if ($mode->endsWithParent && $parent->terminator_end) {
                $mode->terminator_end .= ($mode->end ? "|" : "") . $parent->terminator_end;
            }
        }

        if ($mode->illegal) {
            $mode->illegalRe = $this->langRe($mode->illegal);
        }

        if ($mode->relevance === null) {
            $mode->relevance = 1;
        }

        if (!$mode->contains) {
            $mode->contains = array();
        }

        /** @var Mode[] $expandedContains */
        $expandedContains = array();
        foreach ($mode->contains as &$c) {
            if ($c instanceof \stdClass) {
                Mode::_normalize($c);
            }

            $expandedContains = array_merge($expandedContains, $this->expandOrCloneMode(
                $c === 'self' ? $mode : $c
            ));
        }
        $mode->contains = $expandedContains;

        /** @var Mode $contain */
        foreach ($mode->contains as $contain) {
            $this->compileMode($contain, $mode);
        }

        if ($mode->starts) {
            $this->compileMode($mode->starts, $parent);
        }

        $terminators = new Terminators($this->case_insensitive);
        $mode->terminators = $terminators->_buildModeRegex($mode);

        Mode::_handleDeprecations($mode);
    }

    /**
     * @param array<string, string>|string $rawKeywords
     * @param bool                         $caseSensitive
     *
     * @return array<string, array<int, string|int>>
     */
    private function compileKeywords($rawKeywords, $caseSensitive)
    {
        /** @var array<string, array<int, string|int>> $compiledKeywords */
        $compiledKeywords = array();

        if (is_string($rawKeywords)) {
            $this->splitAndCompile("keyword", $rawKeywords, $compiledKeywords, $caseSensitive);
        } else {
            foreach ($rawKeywords as $className => $rawKeyword) {
                $this->splitAndCompile($className, $rawKeyword, $compiledKeywords, $caseSensitive);
            }
        }

        return $compiledKeywords;
    }

    /**
     * @param string                                $className
     * @param string                                $str
     * @param array<string, array<int, string|int>> $compiledKeywords
     * @param bool                                  $caseSensitive
     *
     * @return void
     */
    private function splitAndCompile($className, $str, array &$compiledKeywords, $caseSensitive)
    {
        if ($caseSensitive) {
            $str = strtolower($str);
        }

        $keywords = explode(' ', $str);

        foreach ($keywords as $keyword) {
            $pair = explode('|', $keyword);
            $providedScore = isset($pair[1]) ? $pair[1] : null;
            $compiledKeywords[$pair[0]] = array($className, $this->scoreForKeyword($pair[0], $providedScore));
        }
    }

    /**
     * @param string $keyword
     * @param string $providedScore
     *
     * @return int
     */
    private function scoreForKeyword($keyword, $providedScore)
    {
        if ($providedScore) {
            return (int) $providedScore;
        }

        return $this->commonKeyword($keyword) ? 0 : 1;
    }

    /**
     * @param string $word
     *
     * @return bool
     */
    private function commonKeyword($word)
    {
        return in_array(strtolower($word), self::$COMMON_KEYWORDS);
    }

    /**
     * Compile the Language definition.
     *
     * @param bool $safeMode
     *
     * @since 9.17.1.0 The 'safeMode' parameter was added.
     *
     * @return void
     */
    public function compile($safeMode)
    {
        if ($this->compiled) {
            return;
        }

        $jr = new JsonRef();
        $jr->decodeRef($this->mode);

        // self is not valid at the top-level
        if (isset($this->mode->contains) && !in_array("self", $this->mode->contains)) {
            if (!$safeMode) {
                throw new \LogicException("`self` is not supported at the top-level of a language.");
            }
            $this->mode->contains = array_filter($this->mode->contains, function ($mode) {
                return $mode !== "self";
            });
        }

        $this->compileMode($this->mode);
    }

    /**
     * @todo Remove in highlight.php 10.x
     *
     * @deprecated 9.16.0 This method should never have been exposed publicly as part of the API.
     *
     * @param \stdClass|null $e
     *
     * @return void
     */
    public function complete(&$e)
    {
        Mode::_normalize($e);
    }
}
