use serenity::{
    async_trait, 
    client::{
        Client, 
        Context, 
        EventHandler
    },
    framework::standard::{
        macros::{
            command, 
            group,
        },
        Args,
        StandardFramework,
        CommandResult
    },
    model::channel::Message, 
    model::gateway::{
        Activity, 
        Ready
    },
    prelude::TypeMapKey,
    model::prelude::ReactionType,
    utils::Colour
    
};
use serde::Deserialize;
use std::{
    collections::HashMap,
    fs::File, 
    io::BufReader, 
    sync::Arc,

};
use regex::Regex;
use tokio::{
    time::Instant,
    time::Duration,
};  
use rand::{
    Rng,
    thread_rng,
};

// struct and traits///////////////////////////////////////////////////
struct WordCount {
    word: String,
    count: usize,
    regex: Regex,
    created: Instant,
}


struct WordCountDatabase;
impl TypeMapKey for WordCountDatabase {
    type Value = HashMap<u64, WordCount>;
}

trait Stack<T> {
    fn peek(&self) -> Option<T>;
}
impl<T> Stack<T> for Vec<T> where T: Copy {
    fn peek(&self) -> Option<T> {
        match self.len() {
            0 => None,
            n => Some(self[n - 1]),
        }
    }
}


// Static consts  ///////////////////////////////////////////////
const MSG_LIMIT: usize = 2000;
const WORD_COUNT_TIME_LIMIT: u64 = 43200; // 12 hours in seconds

// Common functions /////////////////////////////////////////////

macro_rules! help_roll {
    () => ("roll/calc: Calculates a dice expression\n\t> sui roll 2d6 + 1d4 + 1\n")
}

macro_rules! help_count_start {
    () => ("count-start: Starts counting a word or phrase from people in the server\n\t> sui count-start <string_you_want_to_count>\n\t> sui count-start moe\n\t> sui count-start \"so moe?!\"")

}

macro_rules! help_count_stop {
    () => ("count-stop: Stops counting for count-start\n\t> sui count-stop")
}

macro_rules! help_poll {
    () => ("poll: Creates a poll!\n\t> sui poll <description_of_the_poll>\n\t> sui poll \"Is this moe?\"\n\t> sui poll \"Do you love me?\"") 
}

macro_rules! help_calc {
    () => ("calc: Calculator!\n\t> sui calc 1+2*3+(4+5)\n");
}

macro_rules! help {
    () => (concat!(help_roll!(), help_poll!(), help_calc!(), help_count_start!(), help_count_stop!()));
}

macro_rules! wrap_code {
    ($item:expr) => (concat!("```", $item, "```"))
}
async fn say(ctx: &Context, msg: &Message, display: impl std::fmt::Display)  {
    if let Err(why) = msg.channel_id.say(&ctx.http, display).await {
        println!("Error sending message: {:?}", why);
    }
}

trait TryPush {
    fn try_push(&mut self, src: &str) -> bool;
}

impl TryPush for String {
    fn try_push(&mut self, src: &str) -> bool {
        // I'm assuming that len() and capacity() are O(1)
        // If we really want to optimize this, we can always take in the length of src.
        if self.len() + src.len()  <= self.capacity() {
            self.push_str(src);
            return true;
        }
        return false; 
    }
}

fn args_to_string(mut args: Args) -> String {
    let mut ret = String::with_capacity(128);
    ret.push_str(args.single::<String>().unwrap().as_str());
    for arg in args.iter::<String>() {
        ret.push_str(format!(" {}", arg.unwrap()).as_str());
    }

    return ret;
}


fn args_to_string_no_whitespace(mut args: Args) -> String {
    let mut ret = String::with_capacity(128);
    ret.push_str(args.single::<String>().unwrap().as_str());
    for arg in args.iter::<String>() {
        ret.push_str(format!("{}", arg.unwrap()).as_str());
    }

    return ret;
}

fn operate(lhs: i32, rhs: i32, op: &str) -> Option<i32> {
    match op {
        "*" => return Some(lhs * rhs),
        "/" => {
            if rhs == 0 {
                return None;
            }
            return Some(lhs / rhs);
        },
        "+" => return Some(lhs + rhs),
        "-" => return Some(lhs - rhs),
        _ => return None,
    }

}


fn eval_postfix(postfix: &Vec<String>) -> Option<(i32, String)> {
    let mut stack: Vec<i32> = Vec::new();
    let mut rng = thread_rng();
    let mut info: String = String::new();

    let is_op = |s: &str| {
        return s == "+" || s == "-" || s == "*" || s == "/";
    };

    for s in postfix {
        if is_op(s.as_str()) {
            // If we are an operator (O), pop stack twice and get A and B to do BOA.
            let b: i32;
            let a: i32;

            match stack.pop() {
                Some(value) => a = value,
                None => return None,
            }

            match stack.pop() {
                Some(value) => b = value,
                None => return None,
            }

            match operate(b, a, s.as_str()) {
                Some(value) => stack.push(value),
                None => return None,
            }

        }
        else {
            // We will attempt to evaluate the string
            // We accept a raw number, or a number seperated by 'd' (eg. 1d6)

            // Check if string is seperated by 'd'
            let d_values: Vec<&str> = s.split(|c| c =='d' || c == 'D').collect();
            if d_values.len() == 2 {
                let times: i32 = match d_values[0].parse() {
                    Ok(x) => x,
                    Err(_) => return None,
                };

                let sides: i32 = match d_values[1].parse() {
                    Ok(x) => x,
                    Err(_) => return None,   
                };
               

                // Request: they want to see each individual roll
                info.push_str(format!("{}d{} = ", times, sides).as_str());

                let mut sum: i32 = 0;
                if times > 1 {
                    for i in 0..times {
                        println!("{} d {}", i, times);
                        let result = rng.gen_range(1, sides + 1);
                        sum += result;
                        if i == 0 {
                            info.push_str(format!("{}", result).as_str());
                        } else {
                            info.push_str(format!(" + {}", result).as_str());
                        }
                    }
                    info.push_str(format!(" = {}\n", sum).as_str());
                }
                else {
                    sum = rng.gen_range(times, sides + 1);
                    info.push_str(format!("{}\n", sum).as_str());
                }
                stack.push(sum);

            }
            else {
                // otherwise, assume it's a parsable number
                match s.parse::<i32>() {
                    Ok(value) => stack.push(value),
                    Err(_) => return None,
                }
            }

        }

    }

    match stack.len() {
        1 => return Some((stack.peek().unwrap(), info)),
        _ => return None,
    }
}

fn infix_to_postfix(infix: &str) -> Option<Vec<String>> {
    let mut postfix: Vec<String> = Vec::new();
    let mut stack: Vec<char> = Vec::new();
    let mut operand_buffer: String = String::new();

    let precedence = |c: char| -> u32 {
        match c {
            '*' => return 2,
            '/' => return 2,
            '+' => return 1,
            '-' => return 1,
            _ => return 0,
        }
    };

    let is_op = |c: char| -> bool { return precedence(c) > 0; };
    

    let flush_operand_buffer = |postfix: &mut Vec<String>, operand_buffer: &mut String| {
        if operand_buffer.len() > 0 {
            postfix.push(operand_buffer.clone());
            operand_buffer.clear();
        }
    };

    for c in infix.chars() {
        if is_op(c) {
            flush_operand_buffer(&mut postfix, &mut operand_buffer);

            // pop stack until '(' or something of lower precedence
            while !stack.is_empty() {
                let last_op: char = stack.peek().unwrap();
                if precedence(c) > precedence(last_op) {
                    break;
                }
                else {
                    postfix.push(stack.pop().unwrap().to_string());
                }
                
            }
            stack.push(c);
        }
        else if c == '(' { 
            flush_operand_buffer(&mut postfix, &mut operand_buffer);
            stack.push(c);
        }
        else if c == ')' {
            flush_operand_buffer(&mut postfix, &mut operand_buffer);

            // pop stack until '('
            let mut found_open_braces = false;
            while !stack.is_empty() {
                let last_op: char = stack.pop().unwrap();
                if last_op == '(' {
                    found_open_braces = true;
                    break;
                }
                postfix.push(last_op.to_string());

            }

            if !found_open_braces {
                return None;
            }

        }
        else {
            operand_buffer.push(c);
        }
    }
    
    flush_operand_buffer(&mut postfix, &mut operand_buffer);

    while !stack.is_empty() {
        postfix.push(stack.pop().unwrap().to_string());
    }

    return Some(postfix);

}

// Commands /////////////////////////////////////////////////////
#[group]
#[commands(version, help, roll, poll, count_start, count_stop, calc)]
struct General;


const CALC_REPLIES: &'static[&'static str] = &[
    "K-Korokoroko~! 👀\n",
    "*Takes out :abacus: :face_with_monocle:* \n",
    "Let me open my caculator app ^^;\n",
    "Wait, I hadn't had my tea but er...I'll try! ><\n",
    "Er, that's too hard. Let me ask Karu!\n",
    "*frowns, scratches head* fumu fumu...\n",
];

#[command]
async fn calc(ctx: &Context, msg: &Message, args: Args) -> CommandResult {
    if args.len() == 0 {
        say(ctx, msg, wrap_code!(help_calc!())).await;
        return Ok(());
    }

    let get_error = || -> String {
        return format!("Something went wrong. Check that your expression is correctly written and beware of dividing by zero!\n{}", wrap_code!(help_roll!()));
    };


    match infix_to_postfix(args_to_string_no_whitespace(args).as_str()) {
        Some(postfix) => {
            match eval_postfix(&postfix) {
                Some((result, info)) => {
                    // idk if this is cute
                    
                    let reply_action: &str;
                    {
                        let mut rng = thread_rng();
                        reply_action = CALC_REPLIES[rng.gen_range(0, CALC_REPLIES.len())];
                    }

                    let mut reply: String; 
                    {
                        if info.len() > 0 {
                            reply = format!("{}```{}```I got it! It's: **{}**!", reply_action, info, result); 
                            if reply.len() > MSG_LIMIT {
                                reply = format!("{}I got it! It's: **{}**!", reply_action, result);
                            }
                        }
                        else {
                            reply = format!("{}I got it! It's: **{}**!", reply_action, result);
                        }
                    }

                    say(ctx, msg, reply).await;
                    return Ok(());

                },
                None => {
                    say(ctx, msg, get_error()).await;
                    return Ok(());
                }

            }
        },
        None => {
            say(ctx, msg, get_error()).await;
            return Ok(());
        }
    }

}

#[command]
async fn roll(ctx: &Context, msg: &Message, args: Args) -> CommandResult {
    return calc(&ctx, &msg, args).await;
}


#[command]
async fn poll(ctx: &Context, msg: &Message, args: Args) -> CommandResult {
    if args.len() == 0 {
        say(ctx, msg, wrap_code!(help_poll!())).await;
        return Ok(());
    }

    let description = args_to_string(args);    
   
    if let Err(why) = msg.channel_id.send_message(&ctx.http, |m| { m
        .reactions(vec![
                ReactionType::Unicode("👍".to_string()), 
                ReactionType::Unicode("👎".to_string())
            ].into_iter())
        .embed( |e| { e
            .title(format!("Poll Created By: {}", msg.author.name))
            .footer(|f| f.text("React to vote!"))
            .description(description)
            .colour(Colour(0xffffff))
        })
    }).await {
        println!("Error sending message: {:?}", why);
    }

    return Ok(());

}

#[command]
#[aliases("count-start")]
async fn count_start(ctx: &Context, msg: &Message, args: Args) -> CommandResult {
    if args.len() == 0 {
        say(ctx, msg, wrap_code!(help_count_start!())).await;
        return Ok(());
    }

    // Check if entry exists in Word Count
    {
        let rdata= ctx.data.read().await;
        let rwcd = rdata.get::<WordCountDatabase>().unwrap();
    
        if let Some(rwc) = rwcd.get(&msg.channel_id.0) {
            say(ctx, msg, format!("I'm already counting '{}' for this channel! Stop counting with: {}", rwc.word, wrap_code!(help_count_stop!()))).await;
            return Ok(());
        }
    }

    // If we are here, that means the channel_id does not exist, so we need to write into the
    // word counter
    {
        let mut wdata = ctx.data.write().await;
        let wwcd = wdata.get_mut::<WordCountDatabase>().unwrap();
        let word = args_to_string(args); 

        say(ctx, msg, format!("Counting '{}' this channel ! Remember to stop with: {}", word, wrap_code!(help_count_stop!()))).await;
        let regex = Regex::new(format!("(?i){}", word).as_str()).unwrap();
        wwcd.insert(msg.channel_id.0, WordCount{ word, count: 0, regex, created: Instant::now() });
    } 
    println!("Count started!");
    return Ok(());
}


#[command]
#[aliases("count-stop")]
async fn count_stop(ctx: &Context, msg: &Message, args: Args) -> CommandResult {
    if args.len() > 0 {
        say(ctx, msg, wrap_code!(help_count_stop!())).await;
        return Ok(());
    }
    // check if entry exists in WordCountDatabase
    {
        let rdata = ctx.data.read().await;
        let rwcd = rdata.get::<WordCountDatabase>().unwrap();
        if !rwcd.contains_key(&msg.channel_id.0) {
            say(ctx, msg, format!("You didn't ask me to start counting! Use the command below! :sweat_smile: {}", wrap_code!(help_count_start!()))).await;
            return Ok(());
        }
    }

    // If it exists, stop the counter
    let mut data = ctx.data.write().await;
    let wwcd = data.get_mut::<WordCountDatabase>().unwrap();
    let wwc = wwcd.get(&msg.channel_id.0).unwrap();

    // count - 1 because we count the statement that the command started as well
    say(ctx,msg, format!("I counted a total of **{}** '{}'! Have a nice day! :relaxed:", wwc.count - 1, wwc.word)).await;
    wwcd.remove(&msg.channel_id.0);


    return Ok(());
}


#[command]
async fn version(ctx: &Context, msg: &Message) -> CommandResult{
    say(ctx, msg, "I'm SuiBot v1.0.0, written in Rust!!").await;
    return Ok(());
}

#[command]
async fn help(ctx: &Context, msg: &Message) -> CommandResult {
    say(ctx, msg, wrap_code!(help!())).await;
    return Ok(());
}


#[derive(Deserialize)]
struct Config {
    token: String,
    prefix: String,
}

async fn sched_clean_word_count_database(ctx: &Context) {
    // logic for WordCountDatabase cleanup
    let mut remove: bool = false;
    {
        let data = ctx.data.read().await;
        let rwcd = data.get::<WordCountDatabase>().unwrap();
        for (_, v) in rwcd {
            if v.created.elapsed().as_secs() >= WORD_COUNT_TIME_LIMIT {
                remove = true;
                break;
            }
        }
    }

    if remove {
        let mut data = ctx.data.write().await;
        let wwcd = data.get_mut::<WordCountDatabase>().unwrap();
        wwcd.retain(|_, v| v.created.elapsed().as_secs() <= WORD_COUNT_TIME_LIMIT);
    }
    
}

struct Handler; 
#[async_trait] impl EventHandler for Handler {
    async fn message(&self, ctx: Context, msg: Message) {
        // Check if channel is being tracked
        let mut count: usize = 0;
        {
            let rdata = ctx.data.read().await;
            let rwcd = rdata.get::<WordCountDatabase>().unwrap();
        
            if let Some(rwc) = rwcd.get(&msg.channel_id.0) {
                count = rwc.regex.find_iter(&msg.content).count();
            }
        }

        
        if count > 0  {
            let mut wdata = ctx.data.write().await; 
            let wwcd = wdata.get_mut::<WordCountDatabase>().unwrap();
            if let Some(wcd)= wwcd.get_mut(&msg.channel_id.0) {
                wcd.count += count;
            }
        }
        
    }


    async fn ready(&self, ctx: Context, ready: Ready) {
        println!("{} is connected!", ready.user.name);
        ctx.set_activity(Activity::playing("type 'sui help'")).await;

        let ctx = Arc::new(ctx);
        tokio::spawn(async move {
            loop {
                tokio::time::delay_for(Duration::from_secs(60)).await;
                {
                    sched_clean_word_count_database(&ctx).await;
                }

            }
        });
    }
}

#[tokio::main]
async fn main() {
    let mut client: Client;
    {
        let config: Config;
        {
            let file = File::open("config.json").unwrap();
            let reader = BufReader::new(file);
            config = serde_json::from_reader(reader).unwrap();
        }

        let framework = StandardFramework::new()
            .configure(|c| c
                  .with_whitespace(true)
                  .prefix(config.prefix.as_str()))
            .group(&GENERAL_GROUP);
    
        
        client = Client::builder(&config.token)
            .event_handler(Handler)
            .framework(framework)
            .await
            .unwrap();

        let mut data = client.data.write().await;
        // alias database
        {
            data.insert::<WordCountDatabase>(HashMap::new());
        }

      
    }


    if let Err(why) = client.start().await {
        println!("Client error: {:?}", why);
    }
}
