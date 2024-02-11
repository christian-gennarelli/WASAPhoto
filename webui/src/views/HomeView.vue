<script>
    import HeaderTopBar from '../components/HeaderTopBar.vue'
    import Post from '../components/Post.vue'
    export default {
        data() {
            return { 
                user: {
                    ID: localStorage.getItem("ID"),
                    Username: localStorage.getItem("Username"),
                    Name: localStorage.getItem("Name"),
                    Birthdate: localStorage.getItem("Birthdate"),
                    ProfilePic: localStorage.getItem("ProfilePic"),
                },
                posts: [],
                loading: true,
            }
        },
        methods: {
            async getUserStream(){
                await this.$axios.get(
                    '/users/' + this.user.Username + '/stream',
                    { 
                        headers: {
                            'Authorization': this.user.ID,
                        }
                    }
                ).then((res) => {
                    this.posts = res.data 
                }).catch((e) => {
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
            },
        },
        components: {
            HeaderTopBar: HeaderTopBar,
            Post: Post
        },
        created() {
            this.getUserStream();
            this.loading = false
        },
    }
</script>

<template>
    <div v-if="!loading">
        <HeaderTopBar></HeaderTopBar>
        <p style="text-align: center;" v-if="!posts"> No posts to see yet... start following!</p>
        <div v-else class="home-container"> 
            <Post
                v-if="posts"
                v-for="(post, key) in posts" 
                :key="key"
                :post="post" 
                :user="user"
                @add-like="this.posts[key].Likes ? this.posts[key].Likes.unshift(this.user) : this.posts[key].Likes = [this.user]"
                @remove-like="this.posts[key].Likes = this.posts[key].Likes.filter(u => u.Username !== this.user.Username)"
                @remove-post="this.posts.splice(key, 1)"
                @add-comment="(comment) => this.posts[key].Comments ? this.posts[key].Comments.unshift(comment) : this.posts[key].Comments = [comment]"
                @remove-comment="(commentID) => this.posts[key].Comments = this.posts[key].Comments.filter(c => c.CommentID != commentID)"
            ></Post>
        </div>
    </div>
</template>

<style>

.home-container {
    margin: 10px;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
    grid-auto-rows: auto;
}

</style>