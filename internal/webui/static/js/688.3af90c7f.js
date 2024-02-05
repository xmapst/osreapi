"use strict";
(self["webpackChunkvue_antd_pro"] = self["webpackChunkvue_antd_pro"] || []).push([[688],{

/***/ 23688:
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  "default": function() { return /* binding */ _404; }
});

;// CONCATENATED MODULE: ./node_modules/babel-loader/lib/index.js??clonedRuleSet-39.use[1]!./node_modules/@vue/vue-loader-v15/lib/loaders/templateLoader.js??ruleSet[1].rules[3]!./node_modules/@vue/vue-loader-v15/lib/index.js??vue-loader-options!./src/views/404.vue?vue&type=template&id=88b1cca6&
var render = function render() {
  var _vm = this,
    _c = _vm._self._c;
  return _c('a-result', {
    attrs: {
      "status": "404",
      "title": "404",
      "sub-title": "Sorry, the page you visited does not exist."
    },
    scopedSlots: _vm._u([{
      key: "extra",
      fn: function fn() {
        return [_c('a-button', {
          attrs: {
            "type": "primary"
          },
          on: {
            "click": _vm.toHome
          }
        }, [_vm._v(" Back Home ")])];
      },
      proxy: true
    }])
  });
};
var staticRenderFns = [];

;// CONCATENATED MODULE: ./node_modules/thread-loader/dist/cjs.js!./node_modules/babel-loader/lib/index.js??clonedRuleSet-39.use[1]!./node_modules/@vue/vue-loader-v15/lib/index.js??vue-loader-options!./src/views/404.vue?vue&type=script&lang=js&
/* harmony default export */ var _404vue_type_script_lang_js_ = ({
  name: 'Exception404',
  methods: {
    toHome: function toHome() {
      this.$router.push({
        path: '/'
      });
    }
  }
});
;// CONCATENATED MODULE: ./src/views/404.vue?vue&type=script&lang=js&
 /* harmony default export */ var views_404vue_type_script_lang_js_ = (_404vue_type_script_lang_js_); 
// EXTERNAL MODULE: ./node_modules/@vue/vue-loader-v15/lib/runtime/componentNormalizer.js
var componentNormalizer = __webpack_require__(1001);
;// CONCATENATED MODULE: ./src/views/404.vue





/* normalize component */
;
var component = (0,componentNormalizer/* default */.Z)(
  views_404vue_type_script_lang_js_,
  render,
  staticRenderFns,
  false,
  null,
  null,
  null
  
)

/* harmony default export */ var _404 = (component.exports);

/***/ })

}]);