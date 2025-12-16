import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Greet } from "../../../wailsjs/go/app/App"
import {EventsOn, LogPrint} from "../../../wailsjs/runtime"
import * as React from "react"
// 导入新添加的类型定义和后端绑定
import type { XiaohongshuItem } from "@/lib/types"
import { NextPage, Refresh, GetItems, OnItemClick } from "../../../wailsjs/go/xiaohongshu/Xiaohongshu"

export default function HomePage() {
  const [val, setVal] = React.useState("")
  // 新增状态用于存储列表数据
  const [items, setItems] = React.useState<XiaohongshuItem[]>([])

  const greet = async () => {
    let result = await Greet("来自首页的问候")
    LogPrint(result)
    setVal(result)
  }

  // 下一页功能
  const nextPage = async () => {
    try {
      await NextPage()
      // 重新加载数据
      loadItems()
      LogPrint("下一页功能已调用")
    } catch (error) {
      LogPrint("下一页功能调用失败: " + error)
    }
  }

  // 刷新功能
  const refresh = async () => {
    try {
      await Refresh()
      // 重新加载数据
      loadItems()
      LogPrint("刷新功能已调用")
    } catch (error) {
      LogPrint("刷新功能调用失败: " + error)
    }
  }

  // 加载列表数据
  const loadItems = async () => {
    try {
      const result = await GetItems()
      // 将返回的数据转换为 XiaohongshuItem 数组
      const itemsData: XiaohongshuItem[] = result.map((item: any) => ({
        index: item.index,
        title: item.title,
        coverImageUrl: item.coverImageUrl,
        username: item.username,
        avatarUrl: item.avatarUrl
      }))
      setItems(itemsData)
      LogPrint("数据加载成功，共 " + itemsData.length + " 条记录")
    } catch (error) {
      LogPrint("数据加载失败: " + error)
    }
  }

  // 列表项点击处理
  const handleItemClick = async (index: number) => {
    try {
      await OnItemClick(index)
      LogPrint("列表项 " + index + " 被点击")
    } catch (error) {
      LogPrint("列表项点击处理失败: " + error)
    }
  }

    EventsOn("on-mount",()=>{
        loadItems().then();
    })
  // 组件挂载时加载数据
  React.useEffect(() => {

  }, [])

  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <div className="grid auto-rows-min gap-4 ">
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>欢迎使用小红书应用</CardTitle>
            <CardDescription>
              这是应用的主页，您可以在这里开始探索所有功能。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p>点击下面的按钮测试与Go后端的通信：</p>
            <div className="mt-4">
              <Button onClick={greet}>发送问候</Button>
              {val && <p className="mt-2">响应: {val}</p>}
            </div>
            
            {/* 新增的功能按钮 */}
            <div className="mt-6">
              <h3 className="text-lg font-medium">内容浏览功能</h3>
              <div className="flex gap-2 mt-2">
                <Button onClick={nextPage}>下一页</Button>
                <Button onClick={refresh}>刷新</Button>
              </div>
            </div>
            
            {/* 数据列表展示 */}
            <div className="mt-6">
              <h3 className="text-lg font-medium">内容列表</h3>
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-2">
                {items.map((item) => (
                  <div 
                    key={item.index} 
                    className="bg-white rounded-xl overflow-hidden shadow cursor-pointer hover:shadow-lg transition-shadow duration-300 ease-in-out"
                    onClick={() => handleItemClick(item.index)}
                  >
                    {/* 笔记封面图 */}
                    <div className="aspect-square overflow-hidden">
                      <img 
                        src={item.coverImageUrl} 
                        alt={item.title} 
                        className="w-full h-full object-cover"
                      />
                    </div>
                    
                    {/* 标题文本 */}
                    <div className="p-3">
                      <h4 className="font-medium text-gray-900 line-clamp-2 mb-2">
                        {item.title}
                      </h4>
                      
                      {/* 用户信息 */}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center">
                          <img 
                            src={item.avatarUrl} 
                            alt={item.username} 
                            className="w-6 h-6 rounded-full mr-2"
                          />
                          <span className="text-sm text-gray-600 truncate max-w-[80px]">
                            {item.username}
                          </span>
                        </div>
                        
                        {/* 点赞数与收藏数 */}
                        <div className="flex items-center space-x-2">
                          <div className="flex items-center text-gray-500">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                            </svg>
                            <span className="text-xs">23</span>
                          </div>
                          <div className="flex items-center text-gray-500">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
                            </svg>
                            <span className="text-xs">15</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
        
      
      </div>
    </div>
  )
}