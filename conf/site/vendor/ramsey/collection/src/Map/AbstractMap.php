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

use Ramsey\Collection\AbstractArray;
use Ramsey\Collection\Exception\InvalidArgumentException;

use function array_key_exists;
use function array_keys;
use function in_array;

/**
 * This class provides a basic implementation of `MapInterface`, to minimize the
 * effort required to implement this interface.
 */
abstract class AbstractMap extends AbstractArray implements MapInterface
{
    /**
     * Sets the given value to the given offset in the map.
     *
     * @param mixed $offset The offset to set.
     * @param mixed $value The value to set at the given offset.
     *
     * @throws InvalidArgumentException if the offset provided is `null`.
     */
    public function offsetSet($offset, $value): void
    {
        if ($offset === null) {
            throw new InvalidArgumentException(
                'Map elements are key/value pairs; a key must be provided for '
                . 'value ' . $value
            );
        }

        $this->data[$offset] = $value;
    }

    /**
     * Returns `true` if this map contains a mapping for the specified key.
     *
     * @param mixed $key The key to check in the map.
     */
    public function containsKey($key): bool
    {
        return array_key_exists($key, $this->data);
    }

    /**
     * Returns `true` if this map maps one or more keys to the specified value.
     *
     * This performs a strict type check on the value.
     *
     * @param mixed $value The value to check in the map.
     */
    public function containsValue($value): bool
    {
        return in_array($value, $this->data, true);
    }

    /**
     * Return an array of the keys contained in this map.
     *
     * @return mixed[]
     */
    public function keys(): array
    {
        return array_keys($this->data);
    }

    /**
     * Returns the value to which the specified key is mapped, `null` if this
     * map contains no mapping for the key, or (optionally) `$defaultValue` if
     * this map contains no mapping for the key.
     *
     * @param mixed $key The key to return from the map.
     * @param mixed $defaultValue The default value to use if `$key` is not found.
     *
     * @return mixed|null the value or `null` if the key could not be found.
     */
    public function get($key, $defaultValue = null)
    {
        if (!$this->containsKey($key)) {
            return $defaultValue;
        }

        return $this[$key];
    }

    /**
     * Associates the specified value with the specified key in this map.
     *
     * If the map previously contained a mapping for the key, the old value is
     * replaced by the specified value.
     *
     * @param mixed $key The key to put or replace in the map.
     * @param mixed $value The value to store at `$key`.
     *
     * @return mixed|null the previous value associated with key, or `null` if
     *     there was no mapping for `$key`.
     */
    public function put($key, $value)
    {
        $previousValue = $this->get($key);
        $this[$key] = $value;

        return $previousValue;
    }

    /**
     * Associates the specified value with the specified key in this map only if
     * it is not already set.
     *
     * If there is already a value associated with `$key`, this returns that
     * value without replacing it.
     *
     * @param mixed $key The key to put in the map.
     * @param mixed $value The value to store at `$key`.
     *
     * @return mixed|null the previous value associated with key, or `null` if
     *     there was no mapping for `$key`.
     */
    public function putIfAbsent($key, $value)
    {
        $currentValue = $this->get($key);

        if ($currentValue === null) {
            $this[$key] = $value;
        }

        return $currentValue;
    }

    /**
     * Removes the mapping for a key from this map if it is present.
     *
     * @param mixed $key The key to remove from the map.
     *
     * @return mixed|null the previous value associated with key, or `null` if
     *     there was no mapping for `$key`.
     */
    public function remove($key)
    {
        $previousValue = $this->get($key);
        unset($this[$key]);

        return $previousValue;
    }

    /**
     * Removes the entry for the specified key only if it is currently mapped to
     * the specified value.
     *
     * This performs a strict type check on the value.
     *
     * @param mixed $key The key to remove from the map.
     * @param mixed $value The value to match.
     *
     * @return bool true if the value was removed.
     */
    public function removeIf($key, $value): bool
    {
        if ($this->get($key) === $value) {
            unset($this[$key]);

            return true;
        }

        return false;
    }

    /**
     * Replaces the entry for the specified key only if it is currently mapped
     * to some value.
     *
     * @param mixed $key The key to replace.
     * @param mixed $value The value to set at `$key`.
     *
     * @return mixed|null the previous value associated with key, or `null` if
     *     there was no mapping for `$key`.
     */
    public function replace($key, $value)
    {
        $currentValue = $this->get($key);

        if ($this->containsKey($key)) {
            $this[$key] = $value;
        }

        return $currentValue;
    }

    /**
     * Replaces the entry for the specified key only if currently mapped to the
     * specified value.
     *
     * This performs a strict type check on the value.
     *
     * @param mixed $key The key to remove from the map.
     * @param mixed $oldValue The value to match.
     * @param mixed $newValue The value to use as a replacement.
     *
     * @return bool true if the value was replaced.
     */
    public function replaceIf($key, $oldValue, $newValue): bool
    {
        if ($this->get($key) === $oldValue) {
            $this[$key] = $newValue;

            return true;
        }

        return false;
    }
}
