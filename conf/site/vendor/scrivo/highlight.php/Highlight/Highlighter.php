<?php

/* Copyright (c)
 * - 2006-2013, Ivan Sagalaev (maniac@softwaremaniacs.org), highlight.js
 *              (original author)
 * - 2013-2019, Geert Bergman (geert@scrivo.nl), highlight.php
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
 * @api
 *
 * @since 7.5.0.0
 */
class Highlighter
{
    /**
     * @since 9.12.0.4
     */
    const SPAN_END_TAG = "</span>";

    /** @var bool */
    private $safeMode = true;

    // @TODO In v10.x, this value should be static to match highlight.js behavior
    /** @var array<string, mixed> */
    private $options;

    /** @var string */
    private $modeBuffer = "";

    /** @var string */
    private $result = "";

    /** @var Mode|null */
    private $top = null;

    /** @var Language|null */
    private $language = null;

    /** @var int */
    private $relevance = 0;

    /** @var bool */
    private $ignoreIllegals = false;

    /** @var array<string, Mode> */
    private $continuations = array();

    /** @var RegExMatch */
    private $lastMatch;

    /** @var string The current code we are highlighting */
    private $codeToHighlight;

    /** @var array<string, Language> A mapping of a language ID to a Language definition */
    private static $classMap = array();

    /** @var string[] A list of registered language IDs */
    private static $languages = array();

    /** @var array<string, string> A mapping from alias (key) to main language ID (value) */
    private static $aliases = array();

    public function __construct()
    {
        $this->lastMatch = new RegExMatch(array());
        $this->lastMatch->type = "";
        $this->lastMatch->rule = null;

        // @TODO In v10.x, remove the default value for the `languages` value to follow highlight.js behavior
        $this->options = array(
            'classPrefix' => 'hljs-',
            'tabReplace' => null,
            'useBR' => false,
            'languages' => array(
                "xml", "json", "javascript", "css", "php", "http",
            ),
        );

        self::registerLanguages();
    }

    /**
     * Loop through all of the languages in our `languages` folder and automatically register them all.
     *
     * @since 8.3.0.0
     *
     * @return void
     */
    private static function registerLanguages()
    {
        // Languages that take precedence in the classMap array.
        $languagePath = __DIR__ . DIRECTORY_SEPARATOR . "languages" . DIRECTORY_SEPARATOR;
        foreach (array("xml", "django", "javascript", "matlab", "cpp") as $languageId) {
            $filePath = $languagePath . $languageId . ".json";
            if (is_readable($filePath)) {
                self::registerLanguage($languageId, $filePath);
            }
        }

        $d = @dir($languagePath);
        if ($d) {
            while (($entry = $d->read()) !== false) {
                if (substr($entry, -5) === ".json") {
                    $languageId = substr($entry, 0, -5);
                    $filePath = $languagePath . $entry;
                    if (is_readable($filePath)) {
                        self::registerLanguage($languageId, $filePath);
                    }
                }
            }
            $d->close();
        }

        self::$languages = array_keys(self::$classMap);
    }

    /**
     * Register a language definition with the Highlighter's internal language
     * storage. Languages are stored in a static variable, so they'll be available
     * across all instances. You only need to register a language once.
     *
     * @param string $languageId The unique name of a language
     * @param string $filePath   The file path to the language definition
     * @param bool   $overwrite  Overwrite language if it already exists
     *
     * @return Language The object containing the definition for a language's markup
     */
    public static function registerLanguage($languageId, $filePath, $overwrite = false)
    {
        if (!isset(self::$classMap[$languageId]) || $overwrite) {
            $lang = new Language($languageId, $filePath);
            self::$classMap[$languageId] = $lang;

            if ($lang->aliases) {
                foreach ($lang->aliases as $alias) {
                    self::$aliases[$alias] = $languageId;
                }
            }
        }

        return self::$classMap[$languageId];
    }

    /**
     * @param RegEx|null $re
     * @param string     $lexeme
     *
     * @return bool
     */
    private function testRe($re, $lexeme)
    {
        if (!$re) {
            return false;
        }

        $lastIndex = $re->lastIndex;
        $result = $re->exec($lexeme);
        $re->lastIndex = $lastIndex;

        return $result && $result->index === 0;
    }

    /**
     * @param string $value
     *
     * @return RegEx
     */
    private function escapeRe($value)
    {
        return new RegEx(sprintf('/%s/um', preg_quote($value)));
    }

    /**
     * @param Mode   $mode
     * @param string $lexeme
     *
     * @return Mode|null
     */
    private function endOfMode($mode, $lexeme)
    {
        if ($this->testRe($mode->endRe, $lexeme)) {
            while ($mode->endsParent && $mode->parent) {
                $mode = $mode->parent;
            }

            return $mode;
        }

        if ($mode->endsWithParent) {
            return $this->endOfMode($mode->parent, $lexeme);
        }

        return null;
    }

    /**
     * @param Mode       $mode
     * @param RegExMatch $match
     *
     * @return mixed|null
     */
    private function keywordMatch($mode, $match)
    {
        $kwd = $this->language->case_insensitive ? mb_strtolower($match[0], "UTF-8") : $match[0];

        return isset($mode->keywords[$kwd]) ? $mode->keywords[$kwd] : null;
    }

    /**
     * @param string $className
     * @param string $insideSpan
     * @param bool   $leaveOpen
     * @param bool   $noPrefix
     *
     * @return string
     */
    private function buildSpan($className, $insideSpan, $leaveOpen = false, $noPrefix = false)
    {
        if (!$leaveOpen && $insideSpan === '') {
            return '';
        }

        if (!$className) {
            return $insideSpan;
        }

        $classPrefix = $noPrefix ? "" : $this->options['classPrefix'];
        $openSpan = "<span class=\"" . $classPrefix;
        $closeSpan = $leaveOpen ? "" : self::SPAN_END_TAG;

        $openSpan .= $className . "\">";

        return $openSpan . $insideSpan . $closeSpan;
    }

    /**
     * @param string $value
     *
     * @return string
     */
    private function escape($value)
    {
        return htmlspecialchars($value, ENT_NOQUOTES);
    }

    /**
     * @return string
     */
    private function processKeywords()
    {
        if (!$this->top->keywords) {
            return $this->escape($this->modeBuffer);
        }

        $result = "";
        $lastIndex = 0;
        $this->top->lexemesRe->lastIndex = 0;
        $match = $this->top->lexemesRe->exec($this->modeBuffer);

        while ($match) {
            $result .= $this->escape(substr($this->modeBuffer, $lastIndex, $match->index - $lastIndex));
            $keyword_match = $this->keywordMatch($this->top, $match);

            if ($keyword_match) {
                $this->relevance += $keyword_match[1];
                $result .= $this->buildSpan($keyword_match[0], $this->escape($match[0]));
            } else {
                $result .= $this->escape($match[0]);
            }

            $lastIndex = $this->top->lexemesRe->lastIndex;
            $match = $this->top->lexemesRe->exec($this->modeBuffer);
        }

        return $result . $this->escape(substr($this->modeBuffer, $lastIndex));
    }

    /**
     * @return string
     */
    private function processSubLanguage()
    {
        try {
            $hl = new Highlighter();

            // @TODO in v10.x, this should no longer be necessary once `$options` is made static
            $hl->setAutodetectLanguages($this->options['languages']);
            $hl->setClassPrefix($this->options['classPrefix']);
            $hl->setTabReplace($this->options['tabReplace']);

            if (!$this->safeMode) {
                $hl->disableSafeMode();
            }

            $explicit = is_string($this->top->subLanguage);
            if ($explicit && !in_array($this->top->subLanguage, self::$languages)) {
                return $this->escape($this->modeBuffer);
            }

            if ($explicit) {
                $res = $hl->highlight(
                    $this->top->subLanguage,
                    $this->modeBuffer,
                    true,
                    isset($this->continuations[$this->top->subLanguage]) ? $this->continuations[$this->top->subLanguage] : null
                );
            } else {
                $res = $hl->highlightAuto(
                    $this->modeBuffer,
                    count($this->top->subLanguage) ? $this->top->subLanguage : null
                );
            }

            // Counting embedded language score towards the host language may be disabled
            // with zeroing the containing mode relevance. Use case in point is Markdown that
            // allows XML everywhere and makes every XML snippet to have a much larger Markdown
            // score.
            if ($this->top->relevance > 0) {
                $this->relevance += $res->relevance;
            }
            if ($explicit) {
                $this->continuations[$this->top->subLanguage] = $res->top;
            }

            return $this->buildSpan($res->language, $res->value, false, true);
        } catch (\Exception $e) {
            return $this->escape($this->modeBuffer);
        }
    }

    /**
     * @return void
     */
    private function processBuffer()
    {
        if (is_object($this->top) && $this->top->subLanguage) {
            $this->result .= $this->processSubLanguage();
        } else {
            $this->result .= $this->processKeywords();
        }

        $this->modeBuffer = '';
    }

    /**
     * @param Mode $mode
     *
     * @return void
     */
    private function startNewMode($mode)
    {
        $this->result .= $mode->className ? $this->buildSpan($mode->className, "", true) : "";

        $t = clone $mode;
        $t->parent = $this->top;
        $this->top = $t;
    }

    /**
     * @param RegExMatch $match
     *
     * @return int
     */
    private function doBeginMatch($match)
    {
        $lexeme = $match[0];
        $newMode = $match->rule;

        if ($newMode && $newMode->endSameAsBegin) {
            $newMode->endRe = $this->escapeRe($lexeme);
        }

        if ($newMode->skip) {
            $this->modeBuffer .= $lexeme;
        } else {
            if ($newMode->excludeBegin) {
                $this->modeBuffer .= $lexeme;
            }
            $this->processBuffer();
            if (!$newMode->returnBegin && !$newMode->excludeBegin) {
                $this->modeBuffer = $lexeme;
            }
        }
        $this->startNewMode($newMode);

        return $newMode->returnBegin ? 0 : strlen($lexeme);
    }

    /**
     * @param RegExMatch $match
     *
     * @return int|null
     */
    private function doEndMatch($match)
    {
        $lexeme = $match[0];
        $matchPlusRemainder = substr($this->codeToHighlight, $match->index);
        $endMode = $this->endOfMode($this->top, $matchPlusRemainder);

        if (!$endMode) {
            return null;
        }

        $origin = $this->top;
        if ($origin->skip) {
            $this->modeBuffer .= $lexeme;
        } else {
            if (!($origin->returnEnd || $origin->excludeEnd)) {
                $this->modeBuffer .= $lexeme;
            }
            $this->processBuffer();
            if ($origin->excludeEnd) {
                $this->modeBuffer = $lexeme;
            }
        }

        do {
            if ($this->top->className) {
                $this->result .= self::SPAN_END_TAG;
            }
            if (!$this->top->skip && !$this->top->subLanguage) {
                $this->relevance += $this->top->relevance;
            }
            $this->top = $this->top->parent;
        } while ($this->top !== $endMode->parent);

        if ($endMode->starts) {
            if ($endMode->endSameAsBegin) {
                $endMode->starts->endRe = $endMode->endRe;
            }

            $this->startNewMode($endMode->starts);
        }

        return $origin->returnEnd ? 0 : strlen($lexeme);
    }

    /**
     * @param string          $textBeforeMatch
     * @param RegExMatch|null $match
     *
     * @return int
     */
    private function processLexeme($textBeforeMatch, $match = null)
    {
        $lexeme = $match ? $match[0] : null;

        // add non-matched text to the current mode buffer
        $this->modeBuffer .= $textBeforeMatch;

        if ($lexeme === null) {
            $this->processBuffer();

            return 0;
        }

        // we've found a 0 width match and we're stuck, so we need to advance
        // this happens when we have badly behaved rules that have optional matchers to the degree that
        // sometimes they can end up matching nothing at all
        // Ref: https://github.com/highlightjs/highlight.js/issues/2140
        if ($this->lastMatch->type === "begin" && $match->type === "end" && $this->lastMatch->index === $match->index && $lexeme === "") {
            // spit the "skipped" character that our regex choked on back into the output sequence
            $this->modeBuffer .= substr($this->codeToHighlight, $match->index, 1);

            return 1;
        }
        $this->lastMatch = $match;

        if ($match->type === "begin") {
            return $this->doBeginMatch($match);
        } elseif ($match->type === "illegal" && !$this->ignoreIllegals) {
            // illegal match, we do not continue processing
            $_modeRaw = isset($this->top->className) ? $this->top->className : "<unnamed>";

            throw new \UnexpectedValueException("Illegal lexeme \"$lexeme\" for mode \"$_modeRaw\"");
        } elseif ($match->type === "end") {
            $processed = $this->doEndMatch($match);

            if ($processed !== null) {
                return $processed;
            }
        }

        // Why might be find ourselves here?  Only one occasion now.  An end match that was
        // triggered but could not be completed.  When might this happen?  When an `endSameasBegin`
        // rule sets the end rule to a specific match.  Since the overall mode termination rule that's
        // being used to scan the text isn't recompiled that means that any match that LOOKS like
        // the end (but is not, because it is not an exact match to the beginning) will
        // end up here.  A definite end match, but when `doEndMatch` tries to "reapply"
        // the end rule and fails to match, we wind up here, and just silently ignore the end.
        //
        // This causes no real harm other than stopping a few times too many.

        $this->modeBuffer .= $lexeme;

        return strlen($lexeme);
    }

    /**
     * Replace tabs for something more usable.
     *
     * @param string $code
     *
     * @return string
     */
    private function replaceTabs($code)
    {
        if ($this->options['tabReplace'] !== null) {
            return str_replace("\t", $this->options['tabReplace'], $code);
        }

        return $code;
    }

    /**
     * Set the set of languages used for autodetection. When using
     * autodetection the code to highlight will be probed for every language
     * in this set. Limiting this set to only the languages you want to use
     * will greatly improve highlighting speed.
     *
     * @param string[] $set An array of language games to use for autodetection. This defaults
     *                      to a typical set Web development languages.
     *
     * @return void
     */
    public function setAutodetectLanguages(array $set)
    {
        $this->options['languages'] = array_unique($set);
        self::registerLanguages();
    }

    /**
     * Get the tab replacement string.
     *
     * @return string The tab replacement string
     */
    public function getTabReplace()
    {
        return $this->options['tabReplace'];
    }

    /**
     * Set the tab replacement string. This defaults to NULL: no tabs
     * will be replaced.
     *
     * @param string $tabReplace The tab replacement string
     *
     * @return void
     */
    public function setTabReplace($tabReplace)
    {
        $this->options['tabReplace'] = $tabReplace;
    }

    /**
     * Get the class prefix string.
     *
     * @return string The class prefix string
     */
    public function getClassPrefix()
    {
        return $this->options['classPrefix'];
    }

    /**
     * Set the class prefix string.
     *
     * @param string $classPrefix The class prefix string
     *
     * @return void
     */
    public function setClassPrefix($classPrefix)
    {
        $this->options['classPrefix'] = $classPrefix;
    }

    /**
     * @since 9.17.1.0
     *
     * @return void
     */
    public function enableSafeMode()
    {
        $this->safeMode = true;
    }

    /**
     * @since 9.17.1.0
     *
     * @return void
     */
    public function disableSafeMode()
    {
        $this->safeMode = false;
    }

    /**
     * @param string $name
     *
     * @return Language|null
     */
    private function getLanguage($name)
    {
        if (isset(self::$classMap[$name])) {
            return self::$classMap[$name];
        } elseif (isset(self::$aliases[$name]) && isset(self::$classMap[self::$aliases[$name]])) {
            return self::$classMap[self::$aliases[$name]];
        }

        return null;
    }

    /**
     * Determine whether or not a language definition supports auto detection.
     *
     * @param string $name Language name
     *
     * @return bool
     */
    private function autoDetection($name)
    {
        $lang = $this->getLanguage($name);

        return $lang && !$lang->disableAutodetect;
    }

    /**
     * Core highlighting function. Accepts a language name, or an alias, and a
     * string with the code to highlight. Returns an object with the following
     * properties:
     * - relevance (int)
     * - value (an HTML string with highlighting markup).
     *
     * @todo In v10.x, change the return type from \stdClass to HighlightResult
     *
     * @param string    $languageName
     * @param string    $code
     * @param bool      $ignoreIllegals
     * @param Mode|null $continuation
     *
     * @throws \DomainException if the requested language was not in this
     *                          Highlighter's language set
     * @throws \Exception       if an invalid regex was given in a language file
     *
     * @return HighlightResult|\stdClass
     */
    public function highlight($languageName, $code, $ignoreIllegals = true, $continuation = null)
    {
        $this->codeToHighlight = $code;
        $this->language = $this->getLanguage($languageName);

        if ($this->language === null) {
            throw new \DomainException("Unknown language: \"$languageName\"");
        }

        $this->language->compile($this->safeMode);
        $this->top = $continuation ? $continuation : $this->language;
        $this->continuations = array();
        $this->result = "";

        for ($current = $this->top; $current !== $this->language; $current = $current->parent) {
            if ($current->className) {
                $this->result = $this->buildSpan($current->className, '', true) . $this->result;
            }
        }

        $this->modeBuffer = "";
        $this->relevance = 0;
        $this->ignoreIllegals = $ignoreIllegals;

        /** @var HighlightResult $res */
        $res = new \stdClass();
        $res->relevance = 0;
        $res->value = "";
        $res->language = "";
        $res->top = null;
        $res->errorRaised = null;

        try {
            $match = null;
            $count = 0;
            $index = 0;

            while ($this->top) {
                $this->top->terminators->lastIndex = $index;
                $match = $this->top->terminators->exec($this->codeToHighlight);

                if (!$match) {
                    break;
                }

                $count = $this->processLexeme(substr($this->codeToHighlight, $index, $match->index - $index), $match);
                $index = $match->index + $count;
            }

            $this->processLexeme(substr($this->codeToHighlight, $index));

            for ($current = $this->top; isset($current->parent); $current = $current->parent) {
                if ($current->className) {
                    $this->result .= self::SPAN_END_TAG;
                }
            }

            $res->relevance = $this->relevance;
            $res->value = $this->replaceTabs($this->result);
            $res->illegal = false;
            $res->language = $this->language->name;
            $res->top = $this->top;

            return $res;
        } catch (\Exception $e) {
            if (strpos($e->getMessage(), "Illegal") !== false) {
                $res->illegal = true;
                $res->relevance = 0;
                $res->value = $this->escape($this->codeToHighlight);

                return $res;
            } elseif ($this->safeMode) {
                $res->relevance = 0;
                $res->value = $this->escape($this->codeToHighlight);
                $res->language = $languageName;
                $res->top = $this->top;
                $res->errorRaised = $e;

                return $res;
            }

            throw $e;
        }
    }

    /**
     * Highlight the given code by highlighting the given code with each
     * registered language and then finding the match with highest accuracy.
     *
     * @param string        $code
     * @param string[]|null $languageSubset When set to null, this method will attempt to highlight $text with each
     *                                      language. Set this to an array of languages of your choice to limit the
     *                                      amount of languages to try.
     *
     * @throws \Exception       if an invalid regex was given in a language file
     * @throws \DomainException if the attempted language to check does not exist
     *
     * @return HighlightResult|\stdClass
     */
    public function highlightAuto($code, $languageSubset = null)
    {
        /** @var HighlightResult $result */
        $result = new \stdClass();
        $result->relevance = 0;
        $result->value = $this->escape($code);
        $result->language = "";
        $secondBest = clone $result;

        if ($languageSubset === null) {
            $optionsLanguages = $this->options['languages'];

            if (is_array($optionsLanguages) && count($optionsLanguages) > 0) {
                $languageSubset = $optionsLanguages;
            } else {
                $languageSubset = self::$languages;
            }
        }

        foreach ($languageSubset as $name) {
            if ($this->getLanguage($name) === null || !$this->autoDetection($name)) {
                continue;
            }

            $current = $this->highlight($name, $code, false);

            if ($current->relevance > $secondBest->relevance) {
                $secondBest = $current;
            }

            if ($current->relevance > $result->relevance) {
                $secondBest = $result;
                $result = $current;
            }
        }

        if ($secondBest->language) {
            $result->secondBest = $secondBest;
        }

        return $result;
    }

    /**
     * Return a list of all supported languages. Using this list in
     * setAutodetectLanguages will turn on autodetection for all supported
     * languages.
     *
     * @param bool $include_aliases specify whether language aliases
     *                              should be included as well
     *
     * @since 9.12.0.3 The `$include_aliases` parameter was added
     * @since 8.3.0.0
     *
     * @return string[] An array of language names
     */
    public function listLanguages($include_aliases = false)
    {
        if ($include_aliases === true) {
            return array_merge(self::$languages, array_keys(self::$aliases));
        }

        return self::$languages;
    }

    /**
     * Returns list of all available aliases for given language name.
     *
     * @param string $name name or alias of language to look-up
     *
     * @throws \DomainException if the requested language was not in this
     *                          Highlighter's language set
     *
     * @since 9.12.0.3
     *
     * @return string[] An array of all aliases associated with the requested
     *                  language name language. Passed-in name is included as
     *                  well.
     */
    public function getAliasesForLanguage($name)
    {
        $language = self::getLanguage($name);

        if ($language === null) {
            throw new \DomainException("Unknown language: $language");
        }

        if ($language->aliases === null) {
            return array($language->name);
        }

        return array_merge(array($language->name), $language->aliases);
    }
}
