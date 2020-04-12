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

use Ramsey\Collection\Exception\CollectionMismatchException;
use Ramsey\Collection\Exception\InvalidArgumentException;
use Ramsey\Collection\Exception\InvalidSortOrderException;
use Ramsey\Collection\Exception\OutOfBoundsException;
use Ramsey\Collection\Exception\ValueExtractionException;
use Ramsey\Collection\Tool\TypeTrait;
use Ramsey\Collection\Tool\ValueExtractorTrait;
use Ramsey\Collection\Tool\ValueToStringTrait;

use function array_filter;
use function array_map;
use function array_merge;
use function array_search;
use function array_udiff;
use function array_uintersect;
use function current;
use function end;
use function in_array;
use function reset;
use function sprintf;
use function unserialize;
use function usort;

/**
 * This class provides a basic implementation of `CollectionInterface`, to
 * minimize the effort required to implement this interface
 */
abstract class AbstractCollection extends AbstractArray implements CollectionInterface
{
    use TypeTrait;
    use ValueToStringTrait;
    use ValueExtractorTrait;

    /**
     * Ensures that this collection contains the specified element.
     *
     * @param mixed $element The element to add to the collection.
     *
     * @return bool `true` if this collection changed as a result of the call.
     *
     * @throws InvalidArgumentException when the element does not match the
     *     specified type for this collection.
     */
    public function add($element): bool
    {
        $this[] = $element;

        return true;
    }

    /**
     * Returns `true` if this collection contains the specified element.
     *
     * @param mixed $element The element to check whether the collection contains.
     * @param bool $strict Whether to perform a strict type check on the value.
     */
    public function contains($element, bool $strict = true): bool
    {
        return in_array($element, $this->data, $strict);
    }

    /**
     * Sets the given value to the given offset in the array.
     *
     * @param mixed|null $offset The position to set the value in the array, or
     *     `null` to append the value to the array.
     * @param mixed $value The value to set at the given offset.
     *
     * @throws InvalidArgumentException when the value does not match the
     *     specified type for this collection.
     */
    public function offsetSet($offset, $value): void
    {
        if ($this->checkType($this->getType(), $value) === false) {
            throw new InvalidArgumentException(
                'Value must be of type ' . $this->getType() . '; value is '
                . $this->toolValueToString($value)
            );
        }

        if ($offset === null) {
            $this->data[] = $value;
        } else {
            $this->data[$offset] = $value;
        }
    }

    /**
     * Removes a single instance of the specified element from this collection,
     * if it is present.
     *
     * @param mixed $element The element to remove from the collection.
     *
     * @return bool `true` if an element was removed as a result of this call.
     */
    public function remove($element): bool
    {
        if (($position = array_search($element, $this->data, true)) !== false) {
            unset($this->data[$position]);

            return true;
        }

        return false;
    }

    /**
     * Returns the values from given property or method.
     *
     * @param string $propertyOrMethod The property or method name to filter by.
     *
     * @return mixed[]
     *
     * @throws ValueExtractionException if property or method is not defined.
     */
    public function column(string $propertyOrMethod): array
    {
        $temp = [];

        foreach ($this->data as $item) {
            $temp[] = $this->extractValue($item, $propertyOrMethod);
        }

        return $temp;
    }

    /**
     * Returns the first item of the collection.
     *
     * @return mixed
     *
     * @throws OutOfBoundsException when the collection is empty.
     */
    public function first()
    {
        if ($this->isEmpty()) {
            throw new OutOfBoundsException('Can\'t determine first item. Collection is empty');
        }

        reset($this->data);

        return current($this->data);
    }

    /**
     * Returns the last item of the collection.
     *
     * @return mixed
     *
     * @throws OutOfBoundsException when the collection is empty.
     */
    public function last()
    {
        if ($this->isEmpty()) {
            throw new OutOfBoundsException('Can\'t determine last item. Collection is empty');
        }

        $item = end($this->data);
        reset($this->data);

        return $item;
    }

    /**
     * Returns a sorted collection.
     *
     * {@inheritdoc}
     *
     * @param string $propertyOrMethod The property or method to sort by.
     * @param string $order The sort order for the resulting collection (one of
     *     this interface's `SORT_*` constants).
     *
     * @return CollectionInterface<mixed, mixed>
     *
     * @throws InvalidSortOrderException if neither "asc" nor "desc" was given
     *     as the order.
     * @throws ValueExtractionException if property or method is not defined.
     */
    public function sort(string $propertyOrMethod, string $order = self::SORT_ASC): CollectionInterface
    {
        if (!in_array($order, [self::SORT_ASC, self::SORT_DESC], true)) {
            throw new InvalidSortOrderException('Invalid sort order given: ' . $order);
        }

        $collection = clone $this;

        usort($collection->data, function ($a, $b) use ($propertyOrMethod, $order) {
            $aValue = $this->extractValue($a, $propertyOrMethod);
            $bValue = $this->extractValue($b, $propertyOrMethod);

            return ($aValue <=> $bValue) * ($order === self::SORT_DESC ? -1 : 1);
        });

        return $collection;
    }

    /**
     * Returns a filtered collection.
     *
     * {@inheritdoc}
     *
     * @param callable $callback A callable to use for filtering elements.
     *
     * @return CollectionInterface<mixed, mixed>
     */
    public function filter(callable $callback): CollectionInterface
    {
        $collection = clone $this;
        $collection->data = array_merge([], array_filter($collection->data, $callback));

        return $collection;
    }

    /**
     * Returns a collection of matching items.
     *
     * {@inheritdoc}
     *
     * @param string $propertyOrMethod The property or method to evaluate.
     * @param mixed  $value The value to match.
     *
     * @return CollectionInterface<mixed, mixed>
     *
     * @throws ValueExtractionException if property or method is not defined.
     */
    public function where(string $propertyOrMethod, $value): CollectionInterface
    {
        return $this->filter(function ($item) use ($propertyOrMethod, $value) {
            $accessorValue = $this->extractValue($item, $propertyOrMethod);

            return $accessorValue === $value;
        });
    }

    /**
     * Applies a callback to each item of the collection.
     *
     * {@inheritdoc}
     *
     * @param callable $callback A callable to apply to each item of the
     *     collection.
     *
     * @return CollectionInterface<mixed, mixed>
     */
    public function map(callable $callback): CollectionInterface
    {
        $collection = clone $this;
        array_map($callback, $collection->data);

        return $collection;
    }

    /**
     * Create a new collection with divergent items between current and given
     * collection.
     *
     * @param CollectionInterface<mixed, mixed> $other The collection to check for divergent
     *     items.
     *
     * @return CollectionInterface<mixed, mixed>
     *
     * @throws CollectionMismatchException if the given collection is not of the
     *     same type.
     */
    public function diff(CollectionInterface $other): CollectionInterface
    {
        if (!$other instanceof static) {
            throw new CollectionMismatchException('Collection must be of type ' . static::class);
        }

        $comparator = function ($a, $b) {
            return $a === $b ? 0 : -1;
        };

        $diffAtoB = array_udiff($this->data, $other->data, $comparator);
        $diffBtoA = array_udiff($other->data, $this->data, $comparator);

        return new static(array_merge($diffAtoB, $diffBtoA));
    }

    /**
     * Create a new collection with intersecting item between current and given
     * collection.
     *
     * @param CollectionInterface<mixed, mixed> $other The collection to check for
     *     intersecting items.
     *
     * @return CollectionInterface<mixed, mixed>
     *
     * @throws CollectionMismatchException if the given collection is not of the
     *     same type.
     */
    public function intersect(CollectionInterface $other): CollectionInterface
    {
        if (!$other instanceof static) {
            throw new CollectionMismatchException('Collection must be of type ' . static::class);
        }

        $intersect = array_uintersect($this->data, $other->data, function ($a, $b) {
            return $a === $b ? 0 : -1;
        });

        return new static($intersect);
    }

    /**
     * Merge current items and items of given collections into a new one.
     *
     * @param CollectionInterface<mixed, mixed> ...$collections The collections to merge.
     *
     * @return CollectionInterface<mixed, mixed>
     *
     * @throws CollectionMismatchException if any of the given collections are not of the same type.
     */
    public function merge(CollectionInterface ...$collections): CollectionInterface
    {
        $temp = [$this->data];

        foreach ($collections as $index => $collection) {
            if (!$collection instanceof static) {
                throw new CollectionMismatchException(
                    sprintf('Collection with index %d must be of type %s', $index, static::class)
                );
            }

            $temp[] = $collection->toArray();
        }

        return new static(array_merge(...$temp));
    }

    /**
     * @inheritDoc
     */
    public function unserialize($serialized): void
    {
        $this->data = unserialize($serialized, ['allowed_classes' => [$this->getType()]]);
    }
}
