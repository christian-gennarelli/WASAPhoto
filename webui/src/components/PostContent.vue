<script>
    import { getImgUrl } from '../functions/getImgUrl';
    import PopupUserlist from './PopupUserlist.vue';
    export default {
    props: ['post', 'containerClass', 'wrapperClass', 'showLikes', 'username', 'liked'],
    emits: ["change-class", "change-like", "delete-post", 'change-show-likes', 'comment-post'],
    methods: {
        getImgUrl,
    },
    data(){
        return {
            comment: '',
        }
    },
    components: { PopupUserlist }
}

</script>

<template>

    <div :class="containerClass">
        <div :class="wrapperClass">
            <div class="post-header">
                <img v-if="wrapperClass" type="button" @click="this.$emit('change-class')" class="exit" src="@/assets/buttons/close.png">
                <img class="author-img" :src="this.getImgUrl('profile_pics/' + post.Author + '.png')">
                <router-link @click="wrapperClass ? this.$emit('change-class') : null" :to="{ name: 'profile', params: {username: post.Author }}">
                    {{ post.Author }}
                </router-link>
                <span> {{ post.CreationDatetime }}</span>
                <img title="Delete post" @click="this.$emit('delete-post')" class="delete-icon" v-if="post.Author == username && containerClass == 'post-container'" src="@/assets/buttons/x-red.png" >
            </div>
            <div class="post-body">
                <img class="post-image" type=button @click="this.$emit('change-class')" :src="this.getImgUrl(post.Photo)">
                <div class="post-like">
                    <PopupUserlist
                        category="Likes"
                        headertxt="Likes"
                        username=""
                        :show="showLikes"
                        :list="post.Likes ? post.Likes : []"
                        @change-show="this.$emit('change-show-likes')"
                        style="font-size: 18px;"
                    ></PopupUserlist>
                </div>
                <img class="like-icon" v-if="liked" type="button" @click="this.$emit('change-like')" src="@/assets/buttons/liked.png">
                <img class="like-icon" v-else type="button" @click="this.$emit('change-like')" src="@/assets/buttons/unliked.png">
                <div class="comments-title" style="margin-left: 5px" @click="this.$emit('change-class')"> Comments: </div> {{ post.Comments ? post.Comments.length : 0 }}
                <div>
                    <router-link @click="this.$emit('change-class')" :to="{ name: 'profile', params: {username: post.Author }}"> <span class="post-body-title"> {{ post.Author }}: </span> </router-link>
                    <span class="post-body-description"> {{ post.Description }} </span> 
                </div>
            </div> 
            <br>
            <div v-if="this.containerClass == 'post-overlay'">
                <span class="comments-title"> Comments </span>
                <textarea v-model="comment" style="display: block; border-radius: 10px; width: 50%; height: 75px" placeholder="Write a comment!" @keyup.enter="this.$emit('comment-post', comment); this.comment=''"></textarea>
                <div style="display: block" v-for="comment, key in post.Comments" :key="key">
                    <div style="display: grid; grid-template-columns: 3fr 1fr;" class="display: inline-block">
                        <div>
                            <router-link @click="this.$emit('change-class')" :to="{ name: 'profile', params: {username: comment.Author }}"> 
                                <span class="post-body-title"> {{ comment.Author }}</span>:
                            </router-link>
                            <span class="post-body-description"> {{ comment.Body }} </span> 
                        </div>
                        <div>
                            <span style="float: right"> {{ comment.CreationDatetime}} </span>
                            <span><img title="Delete comment" @click="this.$emit('delete-comment', comment.CommentID)" style="border-radius: 15px; width: 24px; height: 24px; float: right" v-if="comment.Author == username"  src="@/assets/buttons/x-red.png"></span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

</template>

<style scoped>

a {
    font-size: 18px;
    text-decoration: none;
    color: black
}

a:hover {
    font-weight: bold;
}

.post-container {
    padding: 1%;
    margin: 1%;
    border: 2px solid black;
    border-radius: 15px;
    background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);
    overflow: hidden;
    overflow-wrap: break-word;
    word-wrap: break-word;
    hyphens: auto;
}

.post-header {
    position: relative;
}

.post-header .author-img {
    margin-bottom: 0.5%;
    border-radius: 15px;
    width: 24px; 
    height: 24px;
}

.post-header span {
    padding-left: 1%
}


.post-body-title {
    font-weight: bold;
    font-size: 15px;
}

.post-body-title:hover {
    font-size: 18px;
}

.post-body-description{
    font-size: 15px;
}

.post-body .post-image {
    display: block; 
    max-width: 336px; 
    max-height: 189px; 
    width: auto; 
    height: auto;
    border-radius: 15px;
    box-shadow: rgba(0, 0, 0, 0.35) 15px 5px 0px;
}

.post-popup .post-image {
    display: block; 
    max-width: 100%; 
    max-height: 200px;
    width: auto; 
    height: auto;
    border-radius: 15px;
    max-width: 683px;
    max-height: 384px;
}

.post-like {
    display: inline-block;
    padding-top: 2%;
    font-size: 18px;
}

.post-overlay {
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    background: black;
    overflow: auto;
    z-index: 1;
    overflow-wrap: break-word;
    word-wrap: break-word;
    hyphens: auto;
}

.post-popup {
    margin-top: 20px;
    margin-bottom: 20px;
    padding: 10px;
    overflow: scroll;
    background: #fff;
    border-radius: 5px;
    width: auto;
    max-width: 55%;
    max-height: 80%;
    position: relative;
    background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);
}

.post-popup::-webkit-scrollbar {
  display: none;
}

.comments-title {
    font-weight: bold;
    font-size: 18px;
    display: inline-block
}

</style>