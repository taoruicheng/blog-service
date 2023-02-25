### 目录结构
~~~
configs：配置文件。
docs：文档集合。
global：全局变量。
internal：内部模块。
    dao：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，例如 MySQL、ElasticSearch 等。
    middleware：HTTP 中间件。
    model：模型层，用于存放 model 对象。
    routers：路由相关逻辑处理。
    service：项目核心业务逻辑。
pkg：项目相关的模块包。
storage：项目生成的临时文件。
scripts：各类构建，安装，分析等操作的脚本。
third_party：第三方的资源工具，例如 Swagger UI。
~~~

* pkg包：如果代码可以导入并在其他项目中使用，则应该位于pkg目录下；如果代码不可重用，或不希望别人重用，则放到/internal目录中。
* internal模块：存放私有应用和库代码。如果一些代码，你不希望在其他应用和库中被导入，可以将这部分代码放在/internal 目录下。
在引入其它项目 internal 下的包时，Go 语言会在编译时报错：
> An import of a path containing the element “internal” is disallowed
if the importing code is outside the tree rooted at the parent of the
"internal" directory.

