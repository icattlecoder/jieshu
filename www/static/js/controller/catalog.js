var jieshuApp = angular.module('jieshuApp', ['ui.bootstrap']);

function DropdownCtrl($scope) {}
jieshuApp.controller('catalogCtl', function($scope, $http) {
    var start = 0;
    $scope.btn_text = "加载更多";
    $scope.books = [];
    var fun = function() {
        $scope.btn_text = "...";
        $http({
            method: 'GET',
            url: '/catalog/search',
            params: {
                "start": start,
                "c": CATALOG,
                "t": TAG
            }
        }).success(function(data, status, headers, config) {
            if (data.books.length == 0) {
                $scope.btn_text = "没有了^︵^";
                return
            }
            $scope.books = $scope.books.concat(data.books);
            start += 1;
            $scope.btn_text = "加载更多";
        }).error(function(data, status, headers, config) {});
    };

    //init data
    fun();

    $scope.search = function(ev) {
        if (ev.which == 13) {
            $scope.books = null;
            $scope.btn_text = "加载书籍...";
            $http({
                method: 'GET',
                url: '/search',
                params: {
                    "keyword": $scope.keyword
                }
            }).success(function(data, status, headers, config) {
                if (data.books.length == 0) {
                    $scope.btn_text = "没有了^︵^";
                    return
                }
                $scope.books = data.books;
                $scope.btn_text = "加载更多";
                start += 1;
            }).error(function(data, status, headers, config) {});
        }
    };
    $scope.more = function() {
        fun();
    };

    $scope.in = function(id, thiz) {

        $http({
            method: 'POST',
            url: '/io/do',
            params: {
                "id": id,
                "io": "in"
            }
        }).success(function(data, status, headers, config) {
            if (data.success) {
                thiz.book.In.push(data.email);
                return;
            } else if (data.needLogin) {
                window.location = data.directUrl;
            }
            console && console.log(data);
        }).error(function(data, status, headers, config) {

            console && console.log(data);
        });
    };

    $scope.out = function(id, thiz) {
        $http({
            method: 'POST',
            url: '/io/do',
            params: {
                "id": id,
                "io": "out"
            }
        }).success(function(data, status, headers, config) {
            if (data.success) {
                thiz.book.Out.push(data.email);
                return;
            } else if (data.needLogin) {
                window.location = data.directUrl;
            }
            console && console.log(data);
        }).error(function(data, status, headers, config) {
            console && console.log(data);
        });
    };
});
