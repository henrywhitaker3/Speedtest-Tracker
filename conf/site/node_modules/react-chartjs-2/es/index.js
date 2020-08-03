"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
Object.defineProperty(exports, "Chart", {
  enumerable: true,
  get: function get() {
    return _chart["default"];
  }
});
exports.defaults = exports.Scatter = exports.Bubble = exports.Polar = exports.Radar = exports.HorizontalBar = exports.Bar = exports.Line = exports.Pie = exports.Doughnut = exports["default"] = void 0;

var _react = _interopRequireDefault(require("react"));

var _propTypes = _interopRequireDefault(require("prop-types"));

var _chart = _interopRequireDefault(require("chart.js"));

var _isEqual = _interopRequireDefault(require("lodash/isEqual"));

var _keyBy = _interopRequireDefault(require("lodash/keyBy"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _extends() { _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; }; return _extends.apply(this, arguments); }

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

function _objectWithoutProperties(source, excluded) { if (source == null) return {}; var target = _objectWithoutPropertiesLoose(source, excluded); var key, i; if (Object.getOwnPropertySymbols) { var sourceSymbolKeys = Object.getOwnPropertySymbols(source); for (i = 0; i < sourceSymbolKeys.length; i++) { key = sourceSymbolKeys[i]; if (excluded.indexOf(key) >= 0) continue; if (!Object.prototype.propertyIsEnumerable.call(source, key)) continue; target[key] = source[key]; } } return target; }

function _objectWithoutPropertiesLoose(source, excluded) { if (source == null) return {}; var target = {}; var sourceKeys = Object.keys(source); var key, i; for (i = 0; i < sourceKeys.length; i++) { key = sourceKeys[i]; if (excluded.indexOf(key) >= 0) continue; target[key] = source[key]; } return target; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function _defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } }

function _createClass(Constructor, protoProps, staticProps) { if (protoProps) _defineProperties(Constructor.prototype, protoProps); if (staticProps) _defineProperties(Constructor, staticProps); return Constructor; }

function _inherits(subClass, superClass) { if (typeof superClass !== "function" && superClass !== null) { throw new TypeError("Super expression must either be null or a function"); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, writable: true, configurable: true } }); if (superClass) _setPrototypeOf(subClass, superClass); }

function _setPrototypeOf(o, p) { _setPrototypeOf = Object.setPrototypeOf || function _setPrototypeOf(o, p) { o.__proto__ = p; return o; }; return _setPrototypeOf(o, p); }

function _createSuper(Derived) { var hasNativeReflectConstruct = _isNativeReflectConstruct(); return function _createSuperInternal() { var Super = _getPrototypeOf(Derived), result; if (hasNativeReflectConstruct) { var NewTarget = _getPrototypeOf(this).constructor; result = Reflect.construct(Super, arguments, NewTarget); } else { result = Super.apply(this, arguments); } return _possibleConstructorReturn(this, result); }; }

function _possibleConstructorReturn(self, call) { if (call && (_typeof(call) === "object" || typeof call === "function")) { return call; } return _assertThisInitialized(self); }

function _assertThisInitialized(self) { if (self === void 0) { throw new ReferenceError("this hasn't been initialised - super() hasn't been called"); } return self; }

function _isNativeReflectConstruct() { if (typeof Reflect === "undefined" || !Reflect.construct) return false; if (Reflect.construct.sham) return false; if (typeof Proxy === "function") return true; try { Date.prototype.toString.call(Reflect.construct(Date, [], function () {})); return true; } catch (e) { return false; } }

function _getPrototypeOf(o) { _getPrototypeOf = Object.setPrototypeOf ? Object.getPrototypeOf : function _getPrototypeOf(o) { return o.__proto__ || Object.getPrototypeOf(o); }; return _getPrototypeOf(o); }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var NODE_ENV = typeof process !== 'undefined' && process.env && process.env.NODE_ENV;

var ChartComponent = /*#__PURE__*/function (_React$Component) {
  _inherits(ChartComponent, _React$Component);

  var _super = _createSuper(ChartComponent);

  function ChartComponent() {
    var _this;

    _classCallCheck(this, ChartComponent);

    _this = _super.call(this);

    _defineProperty(_assertThisInitialized(_this), "handleOnClick", function (event) {
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
    });

    _defineProperty(_assertThisInitialized(_this), "ref", function (element) {
      _this.element = element;
    });

    _this.chartInstance = undefined;
    return _this;
  }

  _createClass(ChartComponent, [{
    key: "componentDidMount",
    value: function componentDidMount() {
      this.renderChart();
    }
  }, {
    key: "componentDidUpdate",
    value: function componentDidUpdate() {
      if (this.props.redraw) {
        this.destroyChart();
        this.renderChart();
        return;
      }

      this.updateChart();
    }
  }, {
    key: "shouldComponentUpdate",
    value: function shouldComponentUpdate(nextProps) {
      var _this$props2 = this.props,
          redraw = _this$props2.redraw,
          type = _this$props2.type,
          options = _this$props2.options,
          plugins = _this$props2.plugins,
          legend = _this$props2.legend,
          height = _this$props2.height,
          width = _this$props2.width;

      if (nextProps.redraw === true) {
        return true;
      }

      if (height !== nextProps.height || width !== nextProps.width) {
        return true;
      }

      if (type !== nextProps.type) {
        return true;
      }

      if (!(0, _isEqual["default"])(legend, nextProps.legend)) {
        return true;
      }

      if (!(0, _isEqual["default"])(options, nextProps.options)) {
        return true;
      }

      var nextData = this.transformDataProp(nextProps);

      if (!(0, _isEqual["default"])(this.shadowDataProp, nextData)) {
        return true;
      }

      return !(0, _isEqual["default"])(plugins, nextProps.plugins);
    }
  }, {
    key: "componentWillUnmount",
    value: function componentWillUnmount() {
      this.destroyChart();
    }
  }, {
    key: "transformDataProp",
    value: function transformDataProp(props) {
      var data = props.data;

      if (typeof data == 'function') {
        var node = this.element;
        return data(node);
      } else {
        return data;
      }
    } // Chart.js directly mutates the data.dataset objects by adding _meta proprerty
    // this makes impossible to compare the current and next data changes
    // therefore we memoize the data prop while sending a fake to Chart.js for mutation.
    // see https://github.com/chartjs/Chart.js/blob/master/src/core/core.controller.js#L615-L617

  }, {
    key: "memoizeDataProps",
    value: function memoizeDataProps() {
      if (!this.props.data) {
        return;
      }

      var data = this.transformDataProp(this.props);
      this.shadowDataProp = _objectSpread(_objectSpread({}, data), {}, {
        datasets: data.datasets && data.datasets.map(function (set) {
          return _objectSpread({}, set);
        })
      });
      this.saveCurrentDatasets(); // to remove the dataset metadata from this chart when the chart is destroyed

      return data;
    }
  }, {
    key: "checkDatasets",
    value: function checkDatasets(datasets) {
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
    }
  }, {
    key: "getCurrentDatasets",
    value: function getCurrentDatasets() {
      return this.chartInstance && this.chartInstance.config.data && this.chartInstance.config.data.datasets || [];
    }
  }, {
    key: "saveCurrentDatasets",
    value: function saveCurrentDatasets() {
      var _this2 = this;

      this.datasets = this.datasets || {};
      var currentDatasets = this.getCurrentDatasets();
      currentDatasets.forEach(function (d) {
        _this2.datasets[_this2.props.datasetKeyProvider(d)] = d;
      });
    }
  }, {
    key: "updateChart",
    value: function updateChart() {
      var _this3 = this;

      var options = this.props.options;
      var data = this.memoizeDataProps(this.props);
      if (!this.chartInstance) return;

      if (options) {
        this.chartInstance.options = _chart["default"].helpers.configMerge(this.chartInstance.options, options);
      } // Pipe datasets to chart instance datasets enabling
      // seamless transitions


      var currentDatasets = this.getCurrentDatasets();
      var nextDatasets = data.datasets || [];
      this.checkDatasets(currentDatasets);
      var currentDatasetsIndexed = (0, _keyBy["default"])(currentDatasets, this.props.datasetKeyProvider); // We can safely replace the dataset array, as long as we retain the _meta property
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
              otherProps = _objectWithoutProperties(next, ["data"]); // Merge properties. Notice a weakness here. If a property is removed
          // from next, it will be retained by current and never disappears.
          // Workaround is to set value to null or undefined in next.


          return _objectSpread(_objectSpread({}, current), otherProps);
        } else {
          return next;
        }
      });

      var datasets = data.datasets,
          rest = _objectWithoutProperties(data, ["datasets"]);

      this.chartInstance.config.data = _objectSpread(_objectSpread({}, this.chartInstance.config.data), rest);
      this.chartInstance.update();
    }
  }, {
    key: "renderChart",
    value: function renderChart() {
      var _this$props3 = this.props,
          options = _this$props3.options,
          legend = _this$props3.legend,
          type = _this$props3.type,
          plugins = _this$props3.plugins;
      var node = this.element;
      var data = this.memoizeDataProps();

      if (typeof legend !== 'undefined' && !(0, _isEqual["default"])(ChartComponent.defaultProps.legend, legend)) {
        options.legend = legend;
      }

      this.chartInstance = new _chart["default"](node, {
        type: type,
        data: data,
        options: options,
        plugins: plugins
      });
    }
  }, {
    key: "destroyChart",
    value: function destroyChart() {
      if (!this.chartInstance) {
        return;
      } // Put all of the datasets that have existed in the chart back on the chart
      // so that the metadata associated with this chart get destroyed.
      // This allows the datasets to be used in another chart. This can happen,
      // for example, in a tabbed UI where the chart gets created each time the
      // tab gets switched to the chart and uses the same data).


      this.saveCurrentDatasets();
      var datasets = Object.values(this.datasets);
      this.chartInstance.config.data.datasets = datasets;
      this.chartInstance.destroy();
    }
  }, {
    key: "render",
    value: function render() {
      var _this$props4 = this.props,
          height = _this$props4.height,
          width = _this$props4.width,
          id = _this$props4.id;
      return /*#__PURE__*/_react["default"].createElement("canvas", {
        ref: this.ref,
        height: height,
        width: width,
        id: id,
        onClick: this.handleOnClick
      });
    }
  }]);

  return ChartComponent;
}(_react["default"].Component);

_defineProperty(ChartComponent, "getLabelAsKey", function (d) {
  return d.label;
});

_defineProperty(ChartComponent, "propTypes", {
  data: _propTypes["default"].oneOfType([_propTypes["default"].object, _propTypes["default"].func]).isRequired,
  getDatasetAtEvent: _propTypes["default"].func,
  getElementAtEvent: _propTypes["default"].func,
  getElementsAtEvent: _propTypes["default"].func,
  height: _propTypes["default"].number,
  legend: _propTypes["default"].object,
  onElementsClick: _propTypes["default"].func,
  options: _propTypes["default"].object,
  plugins: _propTypes["default"].arrayOf(_propTypes["default"].object),
  redraw: _propTypes["default"].bool,
  type: function type(props, propName, componentName) {
    if (!_chart["default"].controllers[props[propName]]) {
      return new Error('Invalid chart type `' + props[propName] + '` supplied to' + ' `' + componentName + '`.');
    }
  },
  width: _propTypes["default"].number,
  datasetKeyProvider: _propTypes["default"].func
});

_defineProperty(ChartComponent, "defaultProps", {
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
});

var _default = ChartComponent;
exports["default"] = _default;

var Doughnut = /*#__PURE__*/function (_React$Component2) {
  _inherits(Doughnut, _React$Component2);

  var _super2 = _createSuper(Doughnut);

  function Doughnut() {
    _classCallCheck(this, Doughnut);

    return _super2.apply(this, arguments);
  }

  _createClass(Doughnut, [{
    key: "render",
    value: function render() {
      var _this4 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref) {
          return _this4.chartInstance = _ref && _ref.chartInstance;
        },
        type: "doughnut"
      }));
    }
  }]);

  return Doughnut;
}(_react["default"].Component);

exports.Doughnut = Doughnut;

var Pie = /*#__PURE__*/function (_React$Component3) {
  _inherits(Pie, _React$Component3);

  var _super3 = _createSuper(Pie);

  function Pie() {
    _classCallCheck(this, Pie);

    return _super3.apply(this, arguments);
  }

  _createClass(Pie, [{
    key: "render",
    value: function render() {
      var _this5 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref2) {
          return _this5.chartInstance = _ref2 && _ref2.chartInstance;
        },
        type: "pie"
      }));
    }
  }]);

  return Pie;
}(_react["default"].Component);

exports.Pie = Pie;

var Line = /*#__PURE__*/function (_React$Component4) {
  _inherits(Line, _React$Component4);

  var _super4 = _createSuper(Line);

  function Line() {
    _classCallCheck(this, Line);

    return _super4.apply(this, arguments);
  }

  _createClass(Line, [{
    key: "render",
    value: function render() {
      var _this6 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref3) {
          return _this6.chartInstance = _ref3 && _ref3.chartInstance;
        },
        type: "line"
      }));
    }
  }]);

  return Line;
}(_react["default"].Component);

exports.Line = Line;

var Bar = /*#__PURE__*/function (_React$Component5) {
  _inherits(Bar, _React$Component5);

  var _super5 = _createSuper(Bar);

  function Bar() {
    _classCallCheck(this, Bar);

    return _super5.apply(this, arguments);
  }

  _createClass(Bar, [{
    key: "render",
    value: function render() {
      var _this7 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref4) {
          return _this7.chartInstance = _ref4 && _ref4.chartInstance;
        },
        type: "bar"
      }));
    }
  }]);

  return Bar;
}(_react["default"].Component);

exports.Bar = Bar;

var HorizontalBar = /*#__PURE__*/function (_React$Component6) {
  _inherits(HorizontalBar, _React$Component6);

  var _super6 = _createSuper(HorizontalBar);

  function HorizontalBar() {
    _classCallCheck(this, HorizontalBar);

    return _super6.apply(this, arguments);
  }

  _createClass(HorizontalBar, [{
    key: "render",
    value: function render() {
      var _this8 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref5) {
          return _this8.chartInstance = _ref5 && _ref5.chartInstance;
        },
        type: "horizontalBar"
      }));
    }
  }]);

  return HorizontalBar;
}(_react["default"].Component);

exports.HorizontalBar = HorizontalBar;

var Radar = /*#__PURE__*/function (_React$Component7) {
  _inherits(Radar, _React$Component7);

  var _super7 = _createSuper(Radar);

  function Radar() {
    _classCallCheck(this, Radar);

    return _super7.apply(this, arguments);
  }

  _createClass(Radar, [{
    key: "render",
    value: function render() {
      var _this9 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref6) {
          return _this9.chartInstance = _ref6 && _ref6.chartInstance;
        },
        type: "radar"
      }));
    }
  }]);

  return Radar;
}(_react["default"].Component);

exports.Radar = Radar;

var Polar = /*#__PURE__*/function (_React$Component8) {
  _inherits(Polar, _React$Component8);

  var _super8 = _createSuper(Polar);

  function Polar() {
    _classCallCheck(this, Polar);

    return _super8.apply(this, arguments);
  }

  _createClass(Polar, [{
    key: "render",
    value: function render() {
      var _this10 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref7) {
          return _this10.chartInstance = _ref7 && _ref7.chartInstance;
        },
        type: "polarArea"
      }));
    }
  }]);

  return Polar;
}(_react["default"].Component);

exports.Polar = Polar;

var Bubble = /*#__PURE__*/function (_React$Component9) {
  _inherits(Bubble, _React$Component9);

  var _super9 = _createSuper(Bubble);

  function Bubble() {
    _classCallCheck(this, Bubble);

    return _super9.apply(this, arguments);
  }

  _createClass(Bubble, [{
    key: "render",
    value: function render() {
      var _this11 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref8) {
          return _this11.chartInstance = _ref8 && _ref8.chartInstance;
        },
        type: "bubble"
      }));
    }
  }]);

  return Bubble;
}(_react["default"].Component);

exports.Bubble = Bubble;

var Scatter = /*#__PURE__*/function (_React$Component10) {
  _inherits(Scatter, _React$Component10);

  var _super10 = _createSuper(Scatter);

  function Scatter() {
    _classCallCheck(this, Scatter);

    return _super10.apply(this, arguments);
  }

  _createClass(Scatter, [{
    key: "render",
    value: function render() {
      var _this12 = this;

      return /*#__PURE__*/_react["default"].createElement(ChartComponent, _extends({}, this.props, {
        ref: function ref(_ref9) {
          return _this12.chartInstance = _ref9 && _ref9.chartInstance;
        },
        type: "scatter"
      }));
    }
  }]);

  return Scatter;
}(_react["default"].Component);

exports.Scatter = Scatter;
var defaults = _chart["default"].defaults;
exports.defaults = defaults;