<?php

namespace Highlight;

/**
 * @internal
 *
 * @since 9.16.0.0
 */
final class Terminators
{
    /** @var bool */
    private $caseInsensitive;

    /** @var array<int, Mode|string> */
    private $matchIndexes = array();

    /** @var RegEx|null */
    private $matcherRe = null;

    /** @var array<int, array<int, Mode|string>> */
    private $regexes = array();

    /** @var int */
    private $matchAt = 1;

    /** @var Mode */
    private $mode;

    /** @var int */
    public $lastIndex = 0;

    /**
     * @param bool $caseInsensitive
     */
    public function __construct($caseInsensitive)
    {
        $this->caseInsensitive = $caseInsensitive;
    }

    /**
     * @internal
     *
     * @param Mode $mode
     *
     * @return self
     */
    public function _buildModeRegex($mode)
    {
        $this->mode = $mode;
        $term = null;

        for ($i = 0; $i < count($mode->contains); ++$i) {
            $re = null;
            $term = $mode->contains[$i];

            if ($term->beginKeywords) {
                $re = "\.?(?:" . $term->begin . ")\.?";
            } else {
                $re = $term->begin;
            }

            $this->addRule($term, $re);
        }

        if ($mode->terminator_end) {
            $this->addRule('end', $mode->terminator_end);
        }

        if ($mode->illegal) {
            $this->addRule('illegal', $mode->illegal);
        }

        /** @var array<int, string> $terminators */
        $terminators = array();
        foreach ($this->regexes as $regex) {
            $terminators[] = $regex[1];
        }

        $this->matcherRe = $this->langRe($this->joinRe($terminators, '|'), true);
        $this->lastIndex = 0;

        return $this;
    }

    /**
     * @param string $s
     *
     * @return RegExMatch|null
     */
    public function exec($s)
    {
        if (count($this->regexes) === 0) {
            return null;
        }

        $this->matcherRe->lastIndex = $this->lastIndex;
        $match = $this->matcherRe->exec($s);
        if (!$match) {
            return null;
        }

        /** @var Mode|string $rule */
        $rule = null;
        for ($i = 0; $i < count($match); ++$i) {
            if ($match[$i] !== null && isset($this->matchIndexes[$i])) {
                $rule = $this->matchIndexes[$i];
                break;
            }
        }

        if (is_string($rule)) {
            $match->type = $rule;
            $match->extra = array($this->mode->illegal, $this->mode->terminator_end);
        } else {
            $match->type = "begin";
            $match->rule = $rule;
        }

        return $match;
    }

    /**
     * @param string $value
     * @param bool   $global
     *
     * @return RegEx
     */
    private function langRe($value, $global = false)
    {
        return RegExUtils::langRe($value, $global, $this->caseInsensitive);
    }

    /**
     * @param Mode|string $rule
     * @param string      $regex
     *
     * @return void
     */
    private function addRule($rule, $regex)
    {
        $this->matchIndexes[$this->matchAt] = $rule;
        $this->regexes[] = array($rule, $regex);
        $this->matchAt += $this->reCountMatchGroups($regex) + 1;
    }

    /**
     * joinRe logically computes regexps.join(separator), but fixes the
     * backreferences so they continue to match.
     *
     * it also places each individual regular expression into it's own
     * match group, keeping track of the sequencing of those match groups
     * is currently an exercise for the caller. :-)
     *
     * @param array<int, string> $regexps
     * @param string             $separator
     *
     * @return string
     */
    private function joinRe($regexps, $separator)
    {
        // backreferenceRe matches an open parenthesis or backreference. To avoid
        // an incorrect parse, it additionally matches the following:
        // - [...] elements, where the meaning of parentheses and escapes change
        // - other escape sequences, so we do not misparse escape sequences as
        //   interesting elements
        // - non-matching or lookahead parentheses, which do not capture. These
        //   follow the '(' with a '?'.
        $backreferenceRe = '#\[(?:[^\\\\\]]|\\\.)*\]|\(\??|\\\([1-9][0-9]*)|\\\.#';
        $numCaptures = 0;
        $ret = '';

        $strLen = count($regexps);
        for ($i = 0; $i < $strLen; ++$i) {
            ++$numCaptures;
            $offset = $numCaptures;
            $re = $this->reStr($regexps[$i]);

            if ($i > 0) {
                $ret .= $separator;
            }

            $ret .= "(";

            while (strlen($re) > 0) {
                $matches = array();
                $matchFound = preg_match($backreferenceRe, $re, $matches, PREG_OFFSET_CAPTURE);

                if ($matchFound === 0) {
                    $ret .= $re;
                    break;
                }

                // PHP aliases to match the JS naming conventions
                $match = $matches[0];
                $index = $match[1];

                $ret .= substr($re, 0, $index);
                $re = substr($re, $index + strlen($match[0]));

                if (substr($match[0], 0, 1) === '\\' && isset($matches[1])) {
                    // Adjust the backreference.
                    $ret .= "\\" . strval(intval($matches[1][0]) + $offset);
                } else {
                    $ret .= $match[0];
                    if ($match[0] == "(") {
                        ++$numCaptures;
                    }
                }
            }

            $ret .= ")";
        }

        return $ret;
    }

    /**
     * @param RegEx|string $re
     *
     * @return mixed
     */
    private function reStr($re)
    {
        if ($re && isset($re->source)) {
            return $re->source;
        }

        return $re;
    }

    /**
     * @param RegEx|string $re
     *
     * @return int
     */
    private function reCountMatchGroups($re)
    {
        $results = array();
        $escaped = preg_replace('#(?<!\\\)/#um', '\\/', (string) $re);
        preg_match_all("/{$escaped}|/u", '', $results);

        return count($results) - 1;
    }
}
