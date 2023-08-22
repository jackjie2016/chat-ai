import service from '@/utils/request'

export const getChatLoginApi = (data) => {
  return service({
    url: '/chat/login',
    method: 'get',
    data
  })
}

export const postChatReplyApi = (data) => {
  return service({
    url: '/chat/reply',
    method: 'post',
    data
  })
}

