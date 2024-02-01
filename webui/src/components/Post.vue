<script>
    import { getImgUrl } from '../functions/getImgUrl'
import PostContent from './PostContent.vue'
    export default {
    props: ['post', 'user'],
    emits: ["add-like", "remove-like", "remove-post"],
    data() {
        return {
            opacity: 0,
            visibility: 'hidden',
            containerClass: 'post-container',
            wrapperClass: '',
            likePhotoPath: (this.post.Likes ? this.post.Likes.map(user => user.Username).includes(localStorage.getItem('Username')) : false) ? new URL('/src/assets/liked.png', import.meta.url) : new URL('/src/assets/unliked.png', import.meta.url)
        };
    },
    methods: {
        changeLike() {
            if (this.likePhotoPath.pathname == '/src/assets/liked.png') {
                this.$axios.delete(
                    '/users/' + this.post.Author + '/profile/posts/' + this.post.PostID + '/likes/' + this.user.Username,
                    {
                        headers: {
                            'Authorization': this.user.ID,
                        }
                    }
                ).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.$emit("remove-like")
                this.likePhotoPath = new URL('/src/assets/unliked.png', import.meta.url)
            }
            else {
                this.$axios.put(
                    '/users/' + this.post.Author + '/profile/posts/' + this.post.PostID + '/likes/',
                    null, 
                    {
                        headers: {
                            'Authorization': this.user.ID,
                        }
                    }
                ).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.$emit("add-like")
                this.likePhotoPath = new URL('/src/assets/liked.png', import.meta.url)
            }
        },
        deletePost(){
            this.$axios.delete(
                '/users/' + this.post.Author + '/profile/posts/' + this.post.PostID + '/',
                {
                    headers: {
                        'Authorization': this.user.ID,
                    }
                }
                ).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
            this.$emit('remove-post')
        },
        changeClass(){
            if (this.containerClass == 'post-container'){
                this.containerClass = 'post-overlay'
                this.wrapperClass = 'post-popup'
            } else {
                this.containerClass = 'post-container'
                this.wrapperClass = ''
            }
        },
        
        getImgUrl
    },
    components: { PostContent }
}
</script>

<template>  

    <PostContent
        :post="post"
        :containerClass="containerClass"
        :wrapperClass="wrapperClass"
        :likePhotoPath="likePhotoPath"
        :username="user.Username"
        @change-class="changeClass"
        @change-like="changeLike"
        @delete-post="deletePost"
    ></PostContent>

</template>