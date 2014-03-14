var jieshuApp = angular.module('jieshuApp', []);

jieshuApp.controller('indexCtl', function($scope, $http) {
    var start = 0;
    $scope.btn_text = "加载更多";
    $scope.books = [];
    var fun = function(index) {
        $scope.btn_text = "加载书籍...";
        $http({
            method: 'GET',
            url: '/search',
            params: {
                "start": start,
                "index": index
            }
        }).success(function(data, status, headers, config) {
            // console.log($scope.books)
            if (data.books.length == 0) {
                $scope.btn_text = "没有了^︵^";
                return
            }
            for (var i = data.books.length - 1; i >= 0; i--) {
                $scope.books.push(data.books[i])
            };
            // $scope.books = $scope.books.push(data.books);
            start += 1;
            $scope.btn_text = "加载更多";
        }).error(function(data, status, headers, config) {

        });
    };

    //init data
    fun(true);

    $scope.search = function() {
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
    };
    $scope.more = function() {
        fun(false);
    };

    $scope. in = function(id, thiz) {

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
                window.location = "/login?returnUrl=/io&b=" + id
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
            }
            console && console.log(data);
        }).error(function(data, status, headers, config) {
            console && console.log(data);
        });
    };
});
