var jieshuApp = angular.module('jieshuApp', ['ui.bootstrap']);

var TabsDemoCtrl = function ($scope) {
};

jieshuApp.controller('ioCtl', function($scope, $http) {
    var start = 0;
    alert("kjs")
    $scope.btn_text = "确定";
    $scope.io = function() {
        alert("lksj")
        $scope.btn_text = "...";
        $http({
            method: 'POST',
            url: '/io',
            params: {
                "verify": $scope.verify
            }
        }).success(function(data, status, headers, config) {
            if (data.books.length == 0) {
                $scope.btn_text = "确定";
                return
            }
            $scope.books = data.books;
            $scope.btn_text = "加载更多";
        }).error(function(data, status, headers, config) {});
    };
    $scope.more = function() {
        fun();
    };
});
