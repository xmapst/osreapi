"use strict";
(self["webpackChunkvue_antd_pro"] = self["webpackChunkvue_antd_pro"] || []).push([[113,616,884,560,639,704],{

/***/ 92747:
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  "default": function() { return /* binding */ lang_zh_CN; }
});

// EXTERNAL MODULE: ./node_modules/@babel/runtime/helpers/esm/objectSpread2.js
var objectSpread2 = __webpack_require__(23711);
// EXTERNAL MODULE: ./node_modules/ant-design-vue/es/vc-pagination/locale/zh_CN.js
var zh_CN = __webpack_require__(56785);
// EXTERNAL MODULE: ./node_modules/babel-runtime/helpers/extends.js
var helpers_extends = __webpack_require__(88239);
;// CONCATENATED MODULE: ./node_modules/ant-design-vue/es/vc-calendar/src/locale/zh_CN.js
/* harmony default export */ var locale_zh_CN = ({
  today: '今天',
  now: '此刻',
  backToToday: '返回今天',
  ok: '确定',
  timeSelect: '选择时间',
  dateSelect: '选择日期',
  weekSelect: '选择周',
  clear: '清除',
  month: '月',
  year: '年',
  previousMonth: '上个月 (翻页上键)',
  nextMonth: '下个月 (翻页下键)',
  monthSelect: '选择月份',
  yearSelect: '选择年份',
  decadeSelect: '选择年代',
  yearFormat: 'YYYY年',
  dayFormat: 'D日',
  dateFormat: 'YYYY年M月D日',
  dateTimeFormat: 'YYYY年M月D日 HH时mm分ss秒',
  previousYear: '上一年 (Control键加左方向键)',
  nextYear: '下一年 (Control键加右方向键)',
  previousDecade: '上一年代',
  nextDecade: '下一年代',
  previousCentury: '上一世纪',
  nextCentury: '下一世纪'
});
;// CONCATENATED MODULE: ./node_modules/ant-design-vue/es/time-picker/locale/zh_CN.js
var locale = {
  placeholder: '请选择时间'
};

/* harmony default export */ var time_picker_locale_zh_CN = (locale);
;// CONCATENATED MODULE: ./node_modules/ant-design-vue/es/date-picker/locale/zh_CN.js




var zh_CN_locale = {
  lang: (0,helpers_extends["default"])({
    placeholder: '请选择日期',
    rangePlaceholder: ['开始日期', '结束日期']
  }, locale_zh_CN),
  timePickerLocale: (0,helpers_extends["default"])({}, time_picker_locale_zh_CN)
};

// should add whitespace between char in Button
zh_CN_locale.lang.ok = '确 定';

// All settings at:
// https://github.com/ant-design/ant-design/blob/master/components/date-picker/locale/example.json

/* harmony default export */ var date_picker_locale_zh_CN = (zh_CN_locale);
;// CONCATENATED MODULE: ./node_modules/ant-design-vue/es/calendar/locale/zh_CN.js

/* harmony default export */ var calendar_locale_zh_CN = (date_picker_locale_zh_CN);
;// CONCATENATED MODULE: ./node_modules/ant-design-vue/es/locale/zh_CN.js





/* harmony default export */ var es_locale_zh_CN = ({
  locale: 'zh-cn',
  Pagination: zh_CN/* default */.Z,
  DatePicker: date_picker_locale_zh_CN,
  TimePicker: time_picker_locale_zh_CN,
  Calendar: calendar_locale_zh_CN,
  // locales for all comoponents
  global: {
    placeholder: '请选择'
  },
  Table: {
    filterTitle: '筛选',
    filterConfirm: '确定',
    filterReset: '重置',
    selectAll: '全选当页',
    selectInvert: '反选当页',
    sortTitle: '排序',
    expand: '展开行',
    collapse: '关闭行'
  },
  Modal: {
    okText: '确定',
    cancelText: '取消',
    justOkText: '知道了'
  },
  Popconfirm: {
    cancelText: '取消',
    okText: '确定'
  },
  Transfer: {
    searchPlaceholder: '请输入搜索内容',
    itemUnit: '项',
    itemsUnit: '项'
  },
  Upload: {
    uploading: '文件上传中',
    removeFile: '删除文件',
    uploadError: '上传错误',
    previewFile: '预览文件',
    downloadFile: '下载文件'
  },
  Empty: {
    description: '暂无数据'
  },
  Icon: {
    icon: '图标'
  },
  Text: {
    edit: '编辑',
    copy: '复制',
    copied: '复制成功',
    expand: '展开'
  },
  PageHeader: {
    back: '返回'
  }
});
;// CONCATENATED MODULE: ./node_modules/ant-design-vue/es/locale-provider/zh_CN.js


/* harmony default export */ var locale_provider_zh_CN = (es_locale_zh_CN);
// EXTERNAL MODULE: ./node_modules/moment/locale/zh-cn.js
var zh_cn = __webpack_require__(83839);
var zh_cn_default = /*#__PURE__*/__webpack_require__.n(zh_cn);
// EXTERNAL MODULE: ./src/locales/lang/zh-CN/global.js
var global = __webpack_require__(67121);
// EXTERNAL MODULE: ./src/locales/lang/zh-CN/menu.js
var menu = __webpack_require__(18695);
// EXTERNAL MODULE: ./src/locales/lang/zh-CN/pool.js
var pool = __webpack_require__(55096);
// EXTERNAL MODULE: ./src/locales/lang/zh-CN/system.js
var system = __webpack_require__(19240);
// EXTERNAL MODULE: ./src/locales/lang/zh-CN/task.js
var task = __webpack_require__(97980);
;// CONCATENATED MODULE: ./src/locales/lang/zh-CN.js








var components = {
  antLocale: locale_provider_zh_CN,
  momentName: 'zh-cn',
  momentLocale: (zh_cn_default())
};
/* harmony default export */ var lang_zh_CN = ((0,objectSpread2/* default */.Z)((0,objectSpread2/* default */.Z)((0,objectSpread2/* default */.Z)((0,objectSpread2/* default */.Z)((0,objectSpread2/* default */.Z)((0,objectSpread2/* default */.Z)({
  message: '-',
  'layouts.usermenu.dialog.title': '信息',
  'layouts.usermenu.dialog.content': '您确定要注销吗？',
  'layouts.userLayout.title': 'Ant Design 是西湖区最具影响力的 Web 设计规范'
}, components), global["default"]), menu["default"]), pool["default"]), system["default"]), task["default"]));

/***/ }),

/***/ 67121:
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

__webpack_require__.r(__webpack_exports__);
/* harmony default export */ __webpack_exports__["default"] = ({
  submit: '提交',
  save: '保存',
  'submit.ok': '提交成功',
  'save.ok': '保存成功'
});

/***/ }),

/***/ 18695:
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

__webpack_require__.r(__webpack_exports__);
/* harmony default export */ __webpack_exports__["default"] = ({
  'menu.welcome': '欢迎',
  'menu.home': '主页',
  'menu.dashboard': '仪表盘',
  'menu.task': '任务',
  'menu.system': '系统'
});

/***/ }),

/***/ 55096:
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

__webpack_require__.r(__webpack_exports__);
/* harmony default export */ __webpack_exports__["default"] = ({
  'pool.worker': '工人数',
  'pool.total': '总任务',
  'pool.running': '运行中',
  'pool.waiting': '等待中'
});

/***/ }),

/***/ 19240:
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

__webpack_require__.r(__webpack_exports__);
/* harmony default export */ __webpack_exports__["default"] = ({
  'system.runtime': '运行时',
  'system.disk': '磁盘',
  'system.cpu': '处理器',
  'system.ram': '内存'
});

/***/ }),

/***/ 97980:
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

__webpack_require__.r(__webpack_exports__);
/* harmony default export */ __webpack_exports__["default"] = ({
  'task.name': '名称',
  'task.id': 'ID',
  'task.state': '状态',
  'task.code': '退出码',
  'task.st': '开始时间',
  'task.et': '结束时间',
  'task.rt': '剩余时间',
  'task.pause': '挂起',
  'task.resume': '解挂',
  'task.kill': '强杀',
  'task.more': '更多',
  'task.0': '已停止',
  'task.1': '运行中',
  'task.2': '等待中',
  'task.3': '已挂起',
  'task.-997': '人工强杀',
  'task.-998': '超时终止',
  'task.-999': '系统错误',
  'task.detail': '详情',
  'task.log': '日志',
  'task.operate': '操作',
  'task.success': '成功'
});

/***/ })

}]);