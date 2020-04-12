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

use Highlight\Highlighter;
use Symfony\Component\Finder\Finder;

class MarkupTest extends PHPUnit_Framework_TestCase
{
    private $allowedFailures;

    public function setUp()
    {
        $this->allowedFailures = array(
            array('haskell', 'nested-comments'),
            array('http', 'default'),
        );
    }

    public static function markupTestProvider()
    {
        $testData = array();

        $markupTests = new Finder();
        $markupTests
            ->in(__DIR__ . '/markup/')
            ->name('*.txt')
            ->sortByName()
            ->files()
        ;

        $workspace = array();

        foreach ($markupTests as $markupTest) {
            $language = $markupTest->getRelativePath();

            if (!isset($workspace[$language])) {
                $workspace[$language] = array();
            }

            if (strpos($markupTest->getFilename(), '.expect.txt') !== false) {
                $workspace[$language][$markupTest->getBasename('.expect.txt')]['expected'] = $markupTest->getContents();
            } else {
                $workspace[$language][$markupTest->getBasename('.txt')]['raw'] = $markupTest->getContents();
            }
        }

        foreach ($workspace as $language => $tests) {
            foreach ($tests as $name => $definition) {
                $testData[] = array($language, $name, $definition['raw'], $definition['expected']);
            }
        }

        return $testData;
    }

    /**
     * @dataProvider markupTestProvider
     */
    public function testHighlighter($language, $testName, $raw, $expected)
    {
        if (in_array(array($language, $testName), $this->allowedFailures)) {
            $this->markTestSkipped("The $language $testName test is known to fail for unknown reasons...");
        }

        $hl = new Highlighter();
        $actual = $hl->highlight($language, $raw);

        $this->assertEquals($language, $actual->language);
        $this->assertEquals(
            trim($expected),
            trim($actual->value),
            sprintf('The "%s" markup test failed for the "%s" language', $testName, $language)
        );
    }
}
