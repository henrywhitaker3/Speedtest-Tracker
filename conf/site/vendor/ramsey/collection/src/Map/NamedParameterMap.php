<?php

/**
 * This file is part of the ramsey/collection library
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 *
 * @copyright Copyright (c) Ben Ramsey <ben@benramsey.com>
 * @license http://opensource.org/licenses/MIT MIT
 */

declare(strict_types=1);

namespace Ramsey\Collection\Map;

use Ramsey\Collection\Exception\InvalidArgumentException;
use Ramsey\Collection\Tool\TypeTrait;
use Ramsey\Collection\Tool\ValueToStringTrait;

use function array_combine;
use function array_key_exists;
use function is_int;

/**
 * `NamedParameterMap` represents a mapping of values to a set of named keys
 * that may optionally be typed
 */
class NamedParameterMap extends AbstractMap
{
    use TypeTrait;
    use ValueToStringTrait;

    /**
     * Named parameters defined for this map.
     *
     * @var array<mixed, string>
     */
    protected $namedParameters;

    /**
     * Constructs a new `NamedParameterMap`.
     *
     * @param array<mixed, string> $namedParameters The named parameters defined for this map.
     * @param mixed[] $data An initial set of data to set on this map.
     */
    public function __construct(array $namedParameters, array $data = [])
    {
        $this->namedParameters = $this->filterNamedParameters($namedParameters);
        parent::__construct($data);
    }

    /**
     * Returns named parameters set for this `NamedParameterMap`.
     *
     * @return array<mixed, string>
     */
    public function getNamedParameters(): array
    {
        return $this->namedParameters;
    }

    /**
     * Sets the given value to the given offset in the map.
     *
     * @param mixed $offset The offset to set.
     * @param mixed $value The value to set at the given offset.
     *
     * @throws InvalidArgumentException if the offset provided is not a
     *     defined named parameter, or if the value is not of the type defined
     *     for the given named parameter.
     */
    public function offsetSet($offset, $value): void
    {
        if (!array_key_exists($offset, $this->namedParameters)) {
            throw new InvalidArgumentException(
                'Attempting to set value for unconfigured parameter \''
                . $offset . '\''
            );
        }

        if ($this->checkType($this->namedParameters[$offset], $value) === false) {
            throw new InvalidArgumentException(
                'Value for \'' . $offset . '\' must be of type '
                . $this->namedParameters[$offset] . '; value is '
                . $this->toolValueToString($value)
            );
        }

        $this->data[$offset] = $value;
    }

    /**
     * Given an array of named parameters, constructs a proper mapping of
     * named parameters to types.
     *
     * @param array<mixed, string> $namedParameters The named parameters to filter.
     *
     * @return array<mixed, string>
     */
    protected function filterNamedParameters(array $namedParameters): array
    {
        $names = [];
        $types = [];

        foreach ($namedParameters as $key => $value) {
            if (is_int($key)) {
                $names[] = (string) $value;
                $types[] = 'mixed';
            } else {
                $names[] = $key;
                $types[] = (string) $value;
            }
        }

        return array_combine($names, $types) ?: [];
    }
}
