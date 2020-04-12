var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

function _objectWithoutProperties(obj, keys) { var target = {}; for (var i in obj) { if (keys.indexOf(i) >= 0) continue; if (!Object.prototype.hasOwnProperty.call(obj, i)) continue; target[i] = obj[i]; } return target; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _possibleConstructorReturn(self, call) { if (!self) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return call && (typeof call === "object" || typeof call === "function") ? call : self; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function, not " + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

import React from 'react';
import PropTypes from 'prop-types';
import Chart from 'chart.js';
import isEqual from 'lodash/isEqual';
import keyBy from 'lodash/keyBy';

var NODE_ENV = typeof process !== 'undefined' && process.env && process.env.NODE_ENV;

var ChartComponent = function (_React$Component) {
  _inherits(ChartComponent, _React$Component);

  function ChartComponent() {
    _classCallCheck(this, ChartComponent);

    var _this = _possibleConstructorReturn(this, _React$Component.call(this));

    _this.handleOnClick = function (event) {
      var instance = _this.chartInstance;

      var _this$props = _this.props,
          getDatasetAtEvent = _this$props.getDatasetAtEvent,
          getElementAtEvent = _this$props.getElementAtEvent,
          getElementsAtEvent = _this$props.getElementsAtEvent,
          onElementsClick = _this$props.onElementsClick;


      getDatasetAtEvent && getDatasetAtEvent(instance.getDatasetAtEvent(event), event);
      getElementAtEvent && getElementAtEvent(instance.getElementAtEvent(event), event);
      getElementsAtEvent && getElementsAtEvent(instance.getElementsAtEvent(event), event);
      onElementsClick && onElementsClick(instance.getElementsAtEvent(event), event); // Backward compatibility
    };

    _this.ref = function (element) {
      _this.element = element;
    };

    _this.chartInstance = undefined;
    return _this;
  }

  ChartComponent.prototype.componentDidMount = function componentDidMount() {
    this.renderChart();
  };

  ChartComponent.prototype.componentDidUpdate = function componentDidUpdate() {
    if (this.props.redraw) {
      this.destroyChart();
      this.renderChart();
      return;
    }

    this.updateChart();
  };

  ChartComponent.prototype.shouldComponentUpdate = function shouldComponentUpdate(nextProps) {
    var _props = this.props,
        redraw = _props.redraw,
        type = _props.type,
        options = _props.options,
        plugins = _props.plugins,
        legend = _props.legend,
        height = _props.height,
        width = _props.width;


    if (nextProps.redraw === true) {
      return true;
    }

    if (height !== nextProps.height || width !== nextProps.width) {
      return true;
    }

    if (type !== nextProps.type) {
      return true;
    }

    if (!isEqual(legend, nextProps.legend)) {
      return true;
    }

    if (!isEqual(options, nextProps.options)) {
      return true;
    }

    var nextData = this.transformDataProp(nextProps);

    if (!isEqual(this.shadowDataProp, nextData)) {
      return true;
    }

    return !isEqual(plugins, nextProps.plugins);
  };

  ChartComponent.prototype.componentWillUnmount = function componentWillUnmount() {
    this.destroyChart();
  };

  ChartComponent.prototype.transformDataProp = function transformDataProp(props) {
    var data = props.data;

    if (typeof data == 'function') {
      var node = this.element;
      return data(node);
    } else {
      return data;
    }
  };

  // Chart.js directly mutates the data.dataset objects by adding _meta proprerty
  // this makes impossible to compare the current and next data changes
  // therefore we memoize the data prop while sending a fake to Chart.js for mutation.
  // see https://github.com/chartjs/Chart.js/blob/master/src/core/core.controller.js#L615-L617


  ChartComponent.prototype.memoizeDataProps = function memoizeDataProps() {
    if (!this.props.data) {
      return;
    }

    var data = this.transformDataProp(this.props);

    this.shadowDataProp = _extends({}, data, {
      datasets: data.datasets && data.datasets.map(function (set) {
        return _extends({}, set);
      })
    });

    this.saveCurrentDatasets(); // to remove the dataset metadata from this chart when the chart is destroyed

    return data;
  };

  ChartComponent.prototype.checkDatasets = function checkDatasets(datasets) {
    var isDev = NODE_ENV !== 'production' && NODE_ENV !== 'prod';
    var usingCustomKeyProvider = this.props.datasetKeyProvider !== ChartComponent.getLabelAsKey;
    var multipleDatasets = datasets.length > 1;

    if (isDev && multipleDatasets && !usingCustomKeyProvider) {
      var shouldWarn = false;
      datasets.forEach(function (dataset) {
        if (!dataset.label) {
          shouldWarn = true;
        }
      });

      if (shouldWarn) {
        console.error('[react-chartjs-2] Warning: Each dataset needs a unique key. By default, the "label" property on each dataset is used. Alternatively, you may provide a "datasetKeyProvider" as a prop that returns a unique key.');
      }
    }
  };

  ChartComponent.prototype.getCurrentDatasets = function getCurrentDatasets() {
    return this.chartInstance && this.chartInstance.config.data && this.chartInstance.config.data.datasets || [];
  };

  ChartComponent.prototype.saveCurrentDatasets = function saveCurrentDatasets() {
    var _this2 = this;

    this.datasets = this.datasets || {};
    var currentDatasets = this.getCurrentDatasets();
    currentDatasets.forEach(function (d) {
      _this2.datasets[_this2.props.datasetKeyProvider(d)] = d;
    });
  };

  ChartComponent.prototype.updateChart = function updateChart() {
    var _this3 = this;

    var options = this.props.options;


    var data = this.memoizeDataProps(this.props);

    if (!this.chartInstance) return;

    if (options) {
      this.chartInstance.options = Chart.helpers.configMerge(this.chartInstance.options, options);
    }

    // Pipe datasets to chart instance datasets enabling
    // seamless transitions
    var currentDatasets = this.getCurrentDatasets();
    var nextDatasets = data.datasets || [];
    this.checkDatasets(currentDatasets);

    var currentDatasetsIndexed = keyBy(currentDatasets, this.props.datasetKeyProvider);

    // We can safely replace the dataset array, as long as we retain the _meta property
    // on each dataset.
    this.chartInstance.config.data.datasets = nextDatasets.map(function (next) {
      var current = currentDatasetsIndexed[_this3.props.datasetKeyProvider(next)];

      if (current && current.type === next.type && next.data) {
        // Be robust to no data. Relevant for other update mechanisms as in chartjs-plugin-streaming.
        // The data array must be edited in place. As chart.js adds listeners to it.
        current.data.splice(next.data.length);
        next.data.forEach(function (point, pid) {
          current.data[pid] = next.data[pid];
        });

        var _data = next.data,
            otherProps = _objectWithoutProperties(next, ['data']);
        // Merge properties. Notice a weakness here. If a property is removed
        // from next, it will be retained by current and never disappears.
        // Workaround is to set value to null or undefined in next.


        return _extends({}, current, otherProps);
      } else {
        return next;
      }
    });

    var datasets = data.datasets,
        rest = _objectWithoutProperties(data, ['datasets']);

    this.chartInstance.config.data = _extends({}, this.chartInstance.config.data, rest);

    this.chartInstance.update();
  };

  ChartComponent.prototype.renderChart = function renderChart() {
    var _props2 = this.props,
        options = _props2.options,
        legend = _props2.legend,
        type = _props2.type,
        plugins = _props2.plugins;

    var node = this.element;
    var data = this.memoizeDataProps();

    if (typeof legend !== 'undefined' && !isEqual(ChartComponent.defaultProps.legend, legend)) {
      options.legend = legend;
    }

    this.chartInstance = new Chart(node, {
      type: type,
      data: data,
      options: options,
      plugins: plugins
    });
  };

  ChartComponent.prototype.destroyChart = function destroyChart() {
    if (!this.chartInstance) {
      return;
    }

    // Put all of the datasets that have existed in the chart back on the chart
    // so that the metadata associated with this chart get destroyed.
    // This allows the datasets to be used in another chart. This can happen,
    // for example, in a tabbed UI where the chart gets created each time the
    // tab gets switched to the chart and uses the same data).
    this.saveCurrentDatasets();
    var datasets = Object.values(this.datasets);
    this.chartInstance.config.data.datasets = datasets;

    this.chartInstance.destroy();
  };

  ChartComponent.prototype.render = function render() {
    var _props3 = this.props,
        height = _props3.height,
        width = _props3.width,
        id = _props3.id;


    return React.createElement('canvas', {
      ref: this.ref,
      height: height,
      width: width,
      id: id,
      onClick: this.handleOnClick
    });
  };

  return ChartComponent;
}(React.Component);

ChartComponent.getLabelAsKey = function (d) {
  return d.label;
};

ChartComponent.propTypes = {
  data: PropTypes.oneOfType([PropTypes.object, PropTypes.func]).isRequired,
  getDatasetAtEvent: PropTypes.func,
  getElementAtEvent: PropTypes.func,
  getElementsAtEvent: PropTypes.func,
  height: PropTypes.number,
  legend: PropTypes.object,
  onElementsClick: PropTypes.func,
  options: PropTypes.object,
  plugins: PropTypes.arrayOf(PropTypes.object),
  redraw: PropTypes.bool,
  type: function type(props, propName, componentName) {
    if (!Chart.controllers[props[propName]]) {
      return new Error('Invalid chart type `' + props[propName] + '` supplied to' + ' `' + componentName + '`.');
    }
  },
  width: PropTypes.number,
  datasetKeyProvider: PropTypes.func
};
ChartComponent.defaultProps = {
  legend: {
    display: true,
    position: 'bottom'
  },
  type: 'doughnut',
  height: 150,
  width: 300,
  redraw: false,
  options: {},
  datasetKeyProvider: ChartComponent.getLabelAsKey
};


export default ChartComponent;

export var Doughnut = function (_React$Component2) {
  _inherits(Doughnut, _React$Component2);

  function Doughnut() {
    _classCallCheck(this, Doughnut);

    return _possibleConstructorReturn(this, _React$Component2.apply(this, arguments));
  }

  Doughnut.prototype.render = function render() {
    var _this5 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref) {
        return _this5.chartInstance = _ref && _ref.chartInstance;
      },
      type: 'doughnut'
    }));
  };

  return Doughnut;
}(React.Component);

export var Pie = function (_React$Component3) {
  _inherits(Pie, _React$Component3);

  function Pie() {
    _classCallCheck(this, Pie);

    return _possibleConstructorReturn(this, _React$Component3.apply(this, arguments));
  }

  Pie.prototype.render = function render() {
    var _this7 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref2) {
        return _this7.chartInstance = _ref2 && _ref2.chartInstance;
      },
      type: 'pie'
    }));
  };

  return Pie;
}(React.Component);

export var Line = function (_React$Component4) {
  _inherits(Line, _React$Component4);

  function Line() {
    _classCallCheck(this, Line);

    return _possibleConstructorReturn(this, _React$Component4.apply(this, arguments));
  }

  Line.prototype.render = function render() {
    var _this9 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref3) {
        return _this9.chartInstance = _ref3 && _ref3.chartInstance;
      },
      type: 'line'
    }));
  };

  return Line;
}(React.Component);

export var Bar = function (_React$Component5) {
  _inherits(Bar, _React$Component5);

  function Bar() {
    _classCallCheck(this, Bar);

    return _possibleConstructorReturn(this, _React$Component5.apply(this, arguments));
  }

  Bar.prototype.render = function render() {
    var _this11 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref4) {
        return _this11.chartInstance = _ref4 && _ref4.chartInstance;
      },
      type: 'bar'
    }));
  };

  return Bar;
}(React.Component);

export var HorizontalBar = function (_React$Component6) {
  _inherits(HorizontalBar, _React$Component6);

  function HorizontalBar() {
    _classCallCheck(this, HorizontalBar);

    return _possibleConstructorReturn(this, _React$Component6.apply(this, arguments));
  }

  HorizontalBar.prototype.render = function render() {
    var _this13 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref5) {
        return _this13.chartInstance = _ref5 && _ref5.chartInstance;
      },
      type: 'horizontalBar'
    }));
  };

  return HorizontalBar;
}(React.Component);

export var Radar = function (_React$Component7) {
  _inherits(Radar, _React$Component7);

  function Radar() {
    _classCallCheck(this, Radar);

    return _possibleConstructorReturn(this, _React$Component7.apply(this, arguments));
  }

  Radar.prototype.render = function render() {
    var _this15 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref6) {
        return _this15.chartInstance = _ref6 && _ref6.chartInstance;
      },
      type: 'radar'
    }));
  };

  return Radar;
}(React.Component);

export var Polar = function (_React$Component8) {
  _inherits(Polar, _React$Component8);

  function Polar() {
    _classCallCheck(this, Polar);

    return _possibleConstructorReturn(this, _React$Component8.apply(this, arguments));
  }

  Polar.prototype.render = function render() {
    var _this17 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref7) {
        return _this17.chartInstance = _ref7 && _ref7.chartInstance;
      },
      type: 'polarArea'
    }));
  };

  return Polar;
}(React.Component);

export var Bubble = function (_React$Component9) {
  _inherits(Bubble, _React$Component9);

  function Bubble() {
    _classCallCheck(this, Bubble);

    return _possibleConstructorReturn(this, _React$Component9.apply(this, arguments));
  }

  Bubble.prototype.render = function render() {
    var _this19 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref8) {
        return _this19.chartInstance = _ref8 && _ref8.chartInstance;
      },
      type: 'bubble'
    }));
  };

  return Bubble;
}(React.Component);

export var Scatter = function (_React$Component10) {
  _inherits(Scatter, _React$Component10);

  function Scatter() {
    _classCallCheck(this, Scatter);

    return _possibleConstructorReturn(this, _React$Component10.apply(this, arguments));
  }

  Scatter.prototype.render = function render() {
    var _this21 = this;

    return React.createElement(ChartComponent, _extends({}, this.props, {
      ref: function ref(_ref9) {
        return _this21.chartInstance = _ref9 && _ref9.chartInstance;
      },
      type: 'scatter'
    }));
  };

  return Scatter;
}(React.Component);

export var defaults = Chart.defaults;
export { Chart };