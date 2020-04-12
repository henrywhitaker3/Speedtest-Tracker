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

/**
 * This class provides a basic implementation of `TypedMapInterface`, to
 * minimize the effort required to implement this interface.
 */
abstract class AbstractTypedMap extends AbstractMap implements TypedMapInterface
{
    use TypeTrait;
    use ValueToStringTrait;

    /**
     * Sets the given value to the given offset in the map.
     *
     * @param mixed $offset The offset to set.
     * @param mixed $value The value to set at the given offset.
     *
     * @throws InvalidArgumentException if the offset or value do not match the
     *     expected types.
     */
    public function offsetSet($offset, $value): void
    {
        if ($this->checkType($this->getKeyType(), $offset) === false) {
            throw new InvalidArgumentException(
                'Key must be of type ' . $this->getKeyType() . '; key is '
                . $this->toolValueToString($offset)
            );
        }

        if ($this->checkType($this->getValueType(), $value) === false) {
            throw new InvalidArgumentException(
                'Value must be of type ' . $this->getValueType() . '; value is '
                . $this->toolValueToString($value)
            );
        }

        parent::offsetSet($offset, $value);
    }
}
