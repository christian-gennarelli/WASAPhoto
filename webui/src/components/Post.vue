<script>
    import { getImgUrl } from '../functions/getImgUrl'
import PostContent from './PostContent.vue'
    export default {
    props: ['post', 'user'],
    emits: ["add-like", "remove-like", "remove-post", "add-comment"],
    data() {
        return {
            opacity: 0,
            visibility: 'hidden',
            containerClass: 'post-container',
            wrapperClass: '',
            showLikes: false,
        };
    },
    computed: {
        liked() {
            return (this.post.Likes ? this.post.Likes.map(user => user.Username).includes(localStorage.getItem('Username')) : false) ? true : false
        }
    },
    methods: {
        changeLike() {
            if (this.liked) {
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
                this.liked = false
            }
            else {
                this.$axios.put(
                    '/users/' + this.post.Author + '/profile/posts/' + this.post.PostID + '/likes/' + this.user.Username,
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
                this.liked = true
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
        commentPost(comment){
            this.$axios.post(
                '/users/' + this.post.Author + '/profile/posts/' + this.post.PostID + '/comments/',
                comment,
                {
                    headers: {
                        'Authorization': this.user.ID,
                        'Content-Type': 'text/plain'
                    }
                },
            ).catch((e) => {
                alert(e.response.data.ErrorCode + " " + e.response.data.Description)
            }).then((res)=>{
                this.$emit('add-comment', res.data)
            })
            
        },
        deleteComment(commentID){
            this.$axios.delete(
                '/users/'+ this.post.Author + '/profile/posts/' + this.post.PostID + '/comments/' + commentID,
                {
                    headers: {
                        'Authorization': this.user.ID,
                    }
                },
            ).catch((e)=>{
                alert(e.response.data.ErrorCode + " " + e.response.data.Description)
            })
            this.$emit('remove-comment', commentID)
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
        :showLikes="showLikes"
        :liked="liked"
        :username="user.Username"
        @change-class="changeClass"
        @change-like="changeLike"
        @change-show-likes="this.showLikes = !this.showLikes; if (this.containerClass == 'post-overlay') {changeClass()}"
        @delete-post="deletePost"
        @comment-post="(comment) => commentPost(comment)"
        @delete-comment="(commentID) => deleteComment(commentID)"
    ></PostContent>

</template>