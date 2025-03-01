#[warn(dead_code)]

use std::fs;
use std::io;

// 获取文件
fn get_file(dir: &str) -> io::Result<Vec<String>> {
    let mut file_list: Vec<String> = Vec::new();
    
    for entry in fs::read_dir(dir)? {
        let entry = entry?;
        let path = entry.path();
        if let Some(ext) = path.extension() {
            if ext.to_string_lossy() == "py" {
                file_list.push(path.to_string_lossy().into_owned());
            }
        }
    }
    Ok(file_list)
}

fn main() {}
