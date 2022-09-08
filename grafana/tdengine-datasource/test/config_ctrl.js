"use strict";

Object.defineProperty(exports, "__esModule", {
    value: true
});

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var GenericDatasourceConfigCtrl =
/** @ngInject */
exports.GenericDatasourceConfigCtrl = function GenericDatasourceConfigCtrl($scope, $injector, $q) {
    _classCallCheck(this, GenericDatasourceConfigCtrl);

    if (!this.current.jsonData.user) {
        this.current.jsonData.user = "root";
    }
    if (!this.current.jsonData.password) {
        this.current.jsonData.password = "taosdata";
    }
};

GenericDatasourceConfigCtrl.templateUrl = 'partials/config.html';
//# sourceMappingURL=config_ctrl.js.map
