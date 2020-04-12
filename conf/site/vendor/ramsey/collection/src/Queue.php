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
use Ramsey\Collection\Exception\NoSuchElementException;
use Ramsey\Collection\Tool\TypeTrait;
use Ramsey\Collection\Tool\ValueToStringTrait;

/**
 * This class provides a basic implementation of `QueueInterface`, to minimize
 * the effort required to implement this interface.
 */
class Queue extends AbstractArray implements QueueInterface
{
    use TypeTrait;
    use ValueToStringTrait;

    /**
     * The type of elements stored in this queue.
     *
     * A queue's type is immutable once it is set. For this reason, this
     * property is set private.
     *
     * @var string
     */
    private $queueType;

    /**
     * The index of the head of the queue.
     *
     * @var int
     */
    protected $index = 0;

    /**
     * Constructs a queue object of the specified type, optionally with the
     * specified data.
     *
     * @param string $queueType The type (FQCN) associated with this queue.
     * @param mixed[] $data The initial items to store in the collection.
     */
    public function __construct(string $queueType, array $data = [])
    {
        $this->queueType = $queueType;
        parent::__construct($data);
    }

    /**
     * Sets the given value to the given offset in the queue.
     *
     * Since arbitrary offsets may not be manipulated in a queue, this method
     * serves only to fulfill the `ArrayAccess` interface requirements. It is
     * invoked by other operations when adding values to the queue.
     *
     * @link http://php.net/manual/en/arrayaccess.offsetset.php ArrayAccess::offsetSet()
     *
     * @param mixed|null $offset The offset is ignored and is treated as `null`.
     * @param mixed $value The value to set at the given offset.
     *
     * @throws InvalidArgumentException when the value does not match the
     *     specified type for this queue.
     */
    public function offsetSet($offset, $value): void
    {
        if ($this->checkType($this->getType(), $value) === false) {
            throw new InvalidArgumentException(
                'Value must be of type ' . $this->getType() . '; value is '
                . $this->toolValueToString($value)
            );
        }

        $this->data[] = $value;
    }

    /**
     * Ensures that this queue contains the specified element.
     *
     * This method differs from `offer()` only in that it throws an exception if
     * it cannot add the element to the queue.
     *
     * @see self::offer()
     *
     * @param mixed $element The element to add to this queue.
     *
     * @return bool `true` if this queue changed as a result of the call.
     *
     * @throws InvalidArgumentException when the element does not match the
     *     specified type for this queue.
     */
    public function add($element): bool
    {
        $this[] = $element;

        return true;
    }

    /**
     * Retrieves, but does not remove, the head of this queue.
     *
     * This method differs from `peek()` only in that it throws an exception if
     * this queue is empty.
     *
     * @see self::peek()
     *
     * @return mixed the head of this queue.
     *
     * @throws NoSuchElementException if this queue is empty.
     */
    public function element()
    {
        if ($this->count() === 0) {
            throw new NoSuchElementException(
                'Can\'t return element from Queue. Queue is empty.'
            );
        }

        return $this[$this->index];
    }

    /**
     * Inserts the specified element into this queue.
     *
     * This method differs from `add()` only in that it does not throw an
     * exception if it cannot add the element to the queue.
     *
     * @see self::add()
     *
     * @param mixed $element The element to add to this queue.
     *
     * @return bool `true` if the element was added to this queue, else `false`.
     */
    public function offer($element): bool
    {
        try {
            return $this->add($element);
        } catch (InvalidArgumentException $e) {
            return false;
        }
    }

    /**
     * Retrieves, but does not remove, the head of this queue, or returns `null`
     * if this queue is empty.
     *
     * @see self::element()
     *
     * @return mixed|null the head of this queue, or `null` if this queue is empty.
     */
    public function peek()
    {
        if ($this->count() === 0) {
            return null;
        }

        return $this[$this->index];
    }

    /**
     * Retrieves and removes the head of this queue, or returns `null`
     * if this queue is empty.
     *
     * @see self::remove()
     *
     * @return mixed|null the head of this queue, or `null` if this queue is empty.
     */
    public function poll()
    {
        if ($this->count() === 0) {
            return null;
        }

        $head = $this[$this->index];

        unset($this[$this->index]);
        $this->index++;

        return $head;
    }

    /**
     * Retrieves and removes the head of this queue.
     *
     * This method differs from `poll()` only in that it throws an exception if
     * this queue is empty.
     *
     * @see self::poll()
     *
     * @return mixed the head of this queue.
     *
     * @throws NoSuchElementException if this queue is empty.
     */
    public function remove()
    {
        if ($this->count() === 0) {
            throw new NoSuchElementException('Can\'t return element from Queue. Queue is empty.');
        }

        $head = $this[$this->index];

        unset($this[$this->index]);
        $this->index++;

        return $head;
    }

    /**
     * Returns the type associated with this queue.
     */
    public function getType(): string
    {
        return $this->queueType;
    }
}
