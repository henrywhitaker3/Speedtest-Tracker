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

namespace Ramsey\Collection;

use Ramsey\Collection\Exception\InvalidArgumentException;

/**
 * This class contains the basic implementation of a collection that does not
 * allow duplicated values (a set), to minimize the effort required to implement
 * this specific type of collection.
 */
abstract class AbstractSet extends AbstractCollection
{
    /**
     * Adds the specified element to this set, if it is not already present.
     *
     * @param mixed $element The element to add to the set.
     *
     * @return bool `true` if this set did not already contain the specified
     *     element.
     *
     * @throws InvalidArgumentException when the element does not match the
     *     specified type for this set.
     */
    public function add($element): bool
    {
        if ($this->contains($element)) {
            return false;
        }

        return parent::add($element);
    }

    /**
     * Sets the given value to the given offset in this set, if it is not
     * already present.
     *
     * @param mixed|null $offset The offset is ignored and is treated as `null`.
     * @param mixed $value The value to set at the given offset.
     *
     * @throws InvalidArgumentException when the value does not match the
     *     specified type for this set.
     */
    public function offsetSet($offset, $value): void
    {
        if ($this->contains($value)) {
            return;
        }

        parent::offsetSet($offset, $value);
    }
}
