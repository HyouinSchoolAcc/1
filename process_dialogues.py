import json
import os
from pathlib import Path

def process_dialogues(folder_path, output_file):
    """
    遍历文件夹中的所有JSON文件，筛选dialogue字段长度>=2的文件，
    整合成一个列表并保存为新的JSON文件
    
    Args:
        folder_path: JSON文件所在的文件夹路径
        output_file: 输出文件名
    """
    sampled_data = []
    processed_count = 0
    selected_count = 0
    
    # 获取文件夹路径
    folder = Path(folder_path)
    
    if not folder.exists():
        print(f"错误：文件夹 {folder_path} 不存在")
        return
    
    # 遍历文件夹中的所有JSON文件
    for json_file in folder.glob("*.json"):
        try:
            # 读取JSON文件
            with open(json_file, 'r', encoding='utf-8') as f:
                data = json.load(f)
            
            processed_count += 1
            
            # 检查是否包含dialogue字段且为列表
            if 'dialogue' in data and isinstance(data['dialogue'], list):
                # 检查dialogue列表长度是否>=2
                if len(data['dialogue']) >= 2:
                    sampled_data.append(data)
                    selected_count += 1
                    print(f"已选择: {json_file.name} (dialogue长度: {len(data['dialogue'])})")
                else:
                    print(f"跳过: {json_file.name} (dialogue长度: {len(data['dialogue'])})")
            else:
                print(f"跳过: {json_file.name} (无dialogue字段或格式不正确)")
                
        except json.JSONDecodeError as e:
            print(f"JSON解析错误 {json_file.name}: {e}")
        except Exception as e:
            print(f"处理文件错误 {json_file.name}: {e}")
    
    # 保存整合后的数据
    try:
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(sampled_data, f, ensure_ascii=False, indent=2)
        
        print(f"\n处理完成!")
        print(f"总处理文件数: {processed_count}")
        print(f"选择的文件数: {selected_count}")
        print(f"结果已保存到: {output_file}")
        
    except Exception as e:
        print(f"保存文件错误: {e}")

if __name__ == "__main__":
    # 设置文件夹路径和输出文件名
    folder_path = "presets/presets_kurisu_CN"
    output_file = "sampled_new.json"
    
    # 执行处理
    process_dialogues(folder_path, output_file) 