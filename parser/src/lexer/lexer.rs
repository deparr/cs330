use std::fmt::Display;

pub enum Token {
}

impl Display for Token {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        return match self {
            _ => write!(f, "todo Token::Display")
        }
    }
}

